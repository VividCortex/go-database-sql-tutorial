---
layout: article
title: Surprises, Antipatterns and Limitations
---

Although `database/sql` is simple once you're accustomed to it, you might be
surprised by the subtlety of use cases it supports. This is common to Go's core
libraries.

Resource Exhaustion
===================

As mentioned throughout this site, if you don't use `database/sql` as intended,
you can certainly cause trouble for yourself, usually by consuming some
resources or preventing them from being reused effectively:

* Opening and closing databases can cause exhaustion of resources.
* Failing to read all rows or use `rows.Close()` reserves connections from the pool.
* Using `Query()` for a statement that doesn't return rows will reserve a connection from the pool.
* Failing to be aware of how [prepared statements](prepared.html) work can lead to a lot of extra database activity.

Large uint64 Values
===================

Here's a surprising error. You can't pass big unsigned integers as parameters to
statements if their high bit is set:

<pre class="prettyprint lang-go">
_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64) // Error
</pre>

This will throw an error. Be careful if you use `uint64` values, as they may
start out small and work without error, but increment over time and start
throwing errors.

Connection State Mismatch
=========================

Some things can change connection state, and that can cause problems for two
reasons:

1. Some connection state, such as whether you're in a transaction, should be
	handled through the Go types instead.
2. You might be assuming that your queries run on a single connection when they
	don't.

For example, setting the current database with a `USE` statement is a typical
thing for many people to do. But in Go, it will affect only the connection that
you run it in. Unless you are in a transaction, other statements that you think
are executed on that connection may actually run on different connections gotten
from the pool, so they won't see the effects of such changes.

Additionally, after you've changed the connection, it'll return to the pool and
potentially pollute the state for some other code. This is one of the reasons
why you should never issue BEGIN or COMMIT statements as SQL commands directly,
too.

Database-Specific Syntax
========================

The `database/sql` API provides an abstraction of a row-oriented database, but
specific databases and drivers can differ in behavior and/or syntax, such as
[prepared statement placeholders](prepared.html).

Multiple Result Sets
====================

The Go driver doesn't support multiple result sets from a single query in any
way, and there doesn't seem to be any plan to do that, although there is [a
feature request](https://github.com/golang/go/issues/5171) for
supporting bulk operations such as bulk copy.

This means, among other things, that a stored procedure that returns multiple
result sets will not work correctly.

Invoking Stored Procedures
==========================

Invoking stored procedures is driver-specific, but in the MySQL driver it can't
be done at present. It might seem that you'd be able to call a simple
procedure that returns a single result set, by executing something like this:

<pre class="prettyprint lang-go">
err := db.QueryRow("CALL mydb.myprocedure").Scan(&amp;result) // Error
</pre>

In fact, this won't work. You'll get the following error: _Error
1312: PROCEDURE mydb.myprocedure can't return a result set in the given
context_. This is because MySQL expects the connection to be set into
multi-statement mode, even for a single result, and the driver doesn't currently
do that (though see [this
issue](https://github.com/go-sql-driver/mysql/issues/66)).

Multiple Statement Support
==========================

The `database/sql` doesn't explicitly have multiple statement support, which means
that the behavior of this is backend dependent:

<pre class="prettyprint lang-go">
_, err := db.Exec("DELETE FROM tbl1; DELETE FROM tbl2") // Error/unpredictable result
</pre>

The server is allowed to interpret this however it wants, which can include
returning an error, executing only the first statement, or executing both.

Similarly, there is no way to batch statements in a transaction. Each statement
in a transaction must be executed serially, and the resources in the results,
such as a Row or Rows, must be scanned or closed so the underlying connection is free
for the next statement to use. This differs from the usual behavior when you're
not working with a transaction. In that scenario, it is perfectly possible to
execute a query, loop over the rows, and within the loop make a query to the
database (which will happen on a new connection):

<pre class="prettyprint lang-go">
rows, err := db.Query("select * from tbl1") // Uses connection 1
for rows.Next() {
	err = rows.Scan(&myvariable)
	// The following line will NOT use connection 1, which is already in-use
	db.Query("select * from tbl2 where id = ?", myvariable)
}
</pre>

But transactions are bound to
just one connection, so this isn't possible with a transaction:

<pre class="prettyprint lang-go">
tx, err := db.Begin()
rows, err := tx.Query("select * from tbl1") // Uses tx's connection
for rows.Next() {
	err = rows.Scan(&myvariable)
	// ERROR! tx's connection is already busy!
	tx.Query("select * from tbl2 where id = ?", myvariable)
}
</pre>

Go doesn't stop you from trying, though. For that reason, you may wind up with a
corrupted connection if you attempt to perform another statement before the
first has released its resources and cleaned up after itself.  This also means
that each statement in a transaction results in a separate set of network
round-trips to the database.

**Previous: [The Connection Pool](connection-pool.html)**
**Next: [Related Reading and Resources](references.html)**
