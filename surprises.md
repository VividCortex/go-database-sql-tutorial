---
layout: page
permalink: /surprises/
title: Surprises, Antipatterns and Limitations
tags:
image:
  feature: abstract-5.jpg
share: false
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
* Failing to use prepared statements can lead to a lot of extra database activity.

Large uint64 Values
===================

Here's a surprising error. You can't pass big unsigned integers as parameters to
statements if their high bit is set:

```go
_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64)
```

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
specific databases and drivers can differ in behavior and/or syntax.  One
example is the syntax for placeholder parameters in prepared statements. For
example, comparing MySQL, PostgreSQL, and Oracle:

	MySQL               PostgreSQL            Oracle
	=====               ==========            ======
	WHERE col = ?       WHERE col = $1        WHERE col = :col
	VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

Multiple Result Sets
====================

The Go driver doesn't support multiple result sets from a single query in any
way, and there doesn't seem to be any plan to do that, although there is [a
feature request](https://code.google.com/p/go/issues/detail?id=5171) for
supporting bulk operations such as bulk copy.

This means, among other things, that a stored procedure that returns multiple
result sets will not work correctly.

Invoking Stored Procedures
==========================

Invoking stored procedures is driver-specific, but in the MySQL driver it can't
be done at present. It might seem that you'd be able to call a simple
procedure that returns a single result set, by executing something like this:

```go
err := db.QueryRow("CALL mydb.myprocedure").Scan(&result)
```

In fact, this won't work. You'll get the following error: _Error
1312: PROCEDURE mydb.myprocedure can't return a result set in the given
context_. This is because MySQL expects the connection to be set into
multi-statement mode, even for a single result, and the driver doesn't currently
do that (though see [this
issue](https://github.com/go-sql-driver/mysql/issues/66)).

Multiple Statement Support
==========================

The `database/sql` doesn't offer multiple-statement support. That means you
can't do something like the following:

```go
_, err := db.Exec("DELETE FROM tbl1; DELETE FROM tbl2")
```

Similarly, there is no way to batch statements in a transaction. Each statement
in a transaction must be executed serially, and the resources that it holds
(such as a Result or Rows) must be closed so the underlying connection is free
for another statement to use. This means that each statement in a transaction
results in a separate set of network round-trips to the database.

{% include toc.md %}
