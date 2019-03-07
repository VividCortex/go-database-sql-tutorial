---
layout: article
title: Modifying Data and Using Transactions
---

Now we're ready to see how to modify data and work with transactions. The
distinction might seem artificial if you're used to programming languages that
use a "statement" object for fetching rows as well as updating data, but in Go,
there's an important reason for the difference.

Statements that Modify Data
===========================

Use `Exec()`, preferably with a prepared statement, to accomplish an `INSERT`,
`UPDATE`, `DELETE`, or another statement that doesn't return rows. The following
example shows how to insert a row and inspect metadata about the operation:

<pre class="prettyprint lang-go">
stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
if err != nil {
	log.Fatal(err)
}
res, err := stmt.Exec("Dolly")
if err != nil {
	log.Fatal(err)
}
lastId, err := res.LastInsertId()
if err != nil {
	log.Fatal(err)
}
rowCnt, err := res.RowsAffected()
if err != nil {
	log.Fatal(err)
}
log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
</pre>

Executing the statement produces a `sql.Result` that gives access to statement
metadata: the last inserted ID and the number of rows affected.

What if you don't care about the result? What if you just want to execute a
statement and check if there were any errors, but ignore the result? Wouldn't
the following two statements do the same thing?

<pre class="prettyprint lang-go">
_, err := db.Exec("DELETE FROM users")  // OK
_, err := db.Query("DELETE FROM users") // BAD
</pre>

The answer is no. They do **not** do the same thing, and **you should never use
`Query()` like this.** The `Query()` will return a `sql.Rows`, which reserves a
database connection until the `sql.Rows` is closed.
Since there might be unread data (e.g. more data rows), the connection can not
be used. In the example above, the connection will *never* be released again.
The garbage collector will eventually close the underlying `net.Conn` for you,
but this might take a long time. Moreover the database/sql package keeps
tracking the connection in its pool, hoping that you release it at some point,
so that the connection can be used again.
This anti-pattern is therefore a good way to run out of resources (too many
connections, for example).

Working with Transactions
=========================

In Go, a transaction is essentially an object that reserves a connection to the
datastore. It lets you do all of the operations we've seen thus far, but
guarantees that they'll be executed on the same connection.

You begin a transaction with a call to `db.Begin()`, and close it with a
`Commit()` or `Rollback()` method on the resulting `Tx` variable. Under the
covers, the `Tx` gets a connection from the pool, and reserves it for use only
with that transaction. The methods on the `Tx` map one-for-one to methods you
can call on the database itself, such as `Query()` and so forth.

Prepared statements that are created in a transaction are bound exclusively to
that transaction. See [prepared statements](prepared.html) for more.

You should not mingle the use of transaction-related functions such as `Begin()`
and `Commit()` with SQL statements such as `BEGIN` and `COMMIT` in your SQL
code. Bad things might result:

* The `Tx` objects could remain open, reserving a connection from the pool and not returning it.
* The state of the database could get out of sync with the state of the Go variables representing it.
* You could believe you're executing queries on a single connection, inside of a transaction, when in reality Go has created several connections for you invisibly and some statements aren't part of the transaction.

While you are working inside a transaction you should be careful not to make
calls to the `db` variable. Make all of your calls to the `Tx` variable that you
created with `db.Begin()`. `db` is not in a transaction, only the `Tx` object is.
If you make further calls to `db.Exec()` or similar, those will happen outside
the scope of your transaction, on other connections.

If you need to work with multiple statements that modify connection state, you
need a `Tx` even if you don't want a transaction per se. For example:

* Creating temporary tables, which are only visible to one connection.
* Setting variables, such as MySQL's `SET @var := somevalue` syntax.
* Changing connection options, such as character sets or timeouts.

If you need to do any of these things, you need to bind your activity to a
single connection, and the only way to do that in Go is to use a `Tx`.

**Previous: [Retrieving Result Sets](retrieving.html)**
**Next: [Using Prepared Statements](prepared.html)**
