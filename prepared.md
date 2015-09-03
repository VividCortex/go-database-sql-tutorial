---
layout: article
title: Using Prepared Statements
---

Prepared statements have all the usual benefits in Go: security, efficiency,
convenience. But the way they're implemented is a little different from what
you might be used to, especially with regards to how they interact with some of
the internals of `database/sql`.

Prepared Statements And Connections
===================================

At the database level, a prepared statement is bound to a single database
connection. The typical flow is that the client sends a SQL statement with
placeholders to the server for preparation, the server responds with a statement
ID, and then the client executes the statement by sending its ID and parameters.

In Go, however, connections are not exposed directly to the user of the
`database/sql` package. You don't prepare a statement on a connection. You
prepare it on a `DB` or a `Tx`. And `database/sql` has some convenience
behaviors such as automatic retries. For these reasons, the underlying
association between prepared statements and connections, which exists at the
driver level, is hidden from your code.

Here's how it works:

1. When you prepare a statement, it's prepared on a connection in the pool.
2. The `Stmt` object remembers which connection was used.
3. When you execute the `Stmt`, it tries to use the connection. If it's not
	available because it's closed or busy doing something else, it gets another
	connection from the pool *and re-prepares the statement with the database on
	another connection.*

Because statements will be re-prepared as needed when their original connection
is busy, it's possible for high-concurrency usage of the database, which may
keep a lot of connections busy, to create a large number of prepared statements.
This can result in apparent leaks of statements, statements being prepared and
re-prepared more often than you think, and even running into server-side limits
on the number of statements.

Avoiding Prepared Statements
============================

Go creates prepared statements for you under the covers. A simple
`db.Query(sql, param1, param2)`, for example, works by preparing the sql, then
executing it with the parameters and finally closing the statement.

Sometimes a prepared statement is not what you want, however. There might be
several reasons for this:

1. The database doesn't support prepared statements. When using the MySQL
	driver, for example, you can connect to MemSQL and Sphinx, because they
	support the MySQL wire protocol. But they don't support the "binary" protocol
	that includes prepared statements, so they can fail in confusing ways.
2. The statements aren't reused enough to make them worthwhile, and security
	issues are handled in other ways, so performance overhead is undesired. An
	example of this can be seen at the
	[VividCortex blog](https://vividcortex.com/blog/2014/11/19/analyzing-prepared-statement-performance-with-vividcortex/).

If you don't want to use a prepared statement, you need to use `fmt.Sprint()` or
similar to assemble the SQL, and pass this as the only argument to `db.Query()`
or `db.QueryRow()`. And your driver needs to support plaintext query execution,
which is added in Go 1.1 via the `Execer` and `Queryer` interfaces,
[documented here](http://golang.org/pkg/database/sql/driver/#Execer).

Prepared Statements in Transactions
===================================

Prepared statements that are created in a `Tx` are bound exclusively to
it, so the earlier cautions about repreparing do not apply. When
you operate on a `Tx` object, your actions map directly to the one and only one
connection underlying it.

This also means that prepared statements created inside a `Tx` can't be used
separately from it. Likewise, prepared statements created on a `DB` can't be
used within a transaction, because they will be bound to a different connection.

To use a prepared statement prepared outside the transaction in a `Tx`, you can use
`Tx.Stmt()`, which will create a new transaction-specific statement from the one
prepared outside the transaction. It does this by taking an existing prepared statement,
setting the connection to that of the transaction and repreparing all statements every
time they are executed. This behavior and its implementation are undesirable and there's
even a TODO in the `database/sql` source code to improve it; we advise against using this.

Caution must be exercised when working with prepared statements in
transactions. Consider the following example:

<pre class="prettyprint lang-go">
tx, err := db.Begin()
if err != nil {
	log.Fatal(err)
}
defer tx.Rollback()
stmt, err := tx.Prepare("INSERT INTO foo VALUES (?)")
if err != nil {
	log.Fatal(err)
}
defer stmt.Close() // danger!
for i := 0; i < 10; i++ {
	_, err = stmt.Exec(i)
	if err != nil {
		log.Fatal(err)
	}
}
err = tx.Commit()
if err != nil {
	log.Fatal(err)
}
// stmt.Close() runs here!
</pre>


Before Go 1.4 closing a `*sql.Tx` released the connection associated with it back into the
pool, but the deferred call to Close on the prepared statement was executed
**after** that has happened, which could lead to concurrent access to the
underlying connection, rendering the connection state inconsistent.
If you use Go 1.4 or older, you should make sure the statement is always closed before the transaction is
committed or rolled back. [This issue](https://github.com/golang/go/issues/4459) was fixed in Go 1.4 by [CR 131650043](https://codereview.appspot.com/131650043).

Parameter Placeholder Syntax
============================

The syntax for placeholder parameters in prepared statements is
database-specific. For example, comparing MySQL, PostgreSQL, and Oracle:

	MySQL               PostgreSQL            Oracle
	=====               ==========            ======
	WHERE col = ?       WHERE col = $1        WHERE col = :col
	VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

**Previous: [Modifying Data and Using Transactions](modifying.html)**
**Next: [Handling Errors](errors.html)**
