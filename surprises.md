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

As mentioned throughtout this site, if you don't use `database/sql` as intended,
you can certainly cause trouble for yourself, usually by consuming some
resources or preventing them from being reused effectively:

* Opening and closing databases can cause exhaustion of resources.
* Failing to use `rows.Close()` or `stmt.Close()` reserves connections from the pool.
* Using `Query()` for a statement that doesn't return rows will reserve a connection from the pool.
* Failing to use prepared statements can lead to a lot of extra database activity.

Large uint64 Values
===================

Here's a surprising error. You can't pass big unsigned integers as parameters to
statements if their high bit is set:

	_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64)

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

{% include toc.md %}
