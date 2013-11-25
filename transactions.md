---
layout: page
permalink: /transactions/
title: Working with Transactions
tags: 
image:
  feature: abstract-5.jpg
share: true
---

You begin a transaction with a call to `db.Begin()`, and close it with a `Commit()` or `Rollback()`
method on the resulting `Tx` variable. Under the covers, the `Tx` gets a connection from the pool,
and reserves it for use only with that transaction. The methods on the `Tx` map one-for-one
to methods you can call on the database itself, such as `Query()` and so forth. The main difference is
that when you call the methods on the `Tx`, they are executed exclusively on the connection associated
with the transaction.

Prepared statements that are created in a transaction are bound exclusively to that transaction, and
can't be used outside of it. Likewise, prepared statements created only on a database handle can't be
used within a transaction.

You should not mingle the use of transaction-related functions such as `Begin()` and `Commit()` with
SQL statements such as `BEGIN` and `COMMIT` in your SQL code. Bad things might result:

1. The `Tx` objects could remain open, reserving a connection from the pool and not returning it.
2. The state of the database could get out of sync with the state of the Go variables representing it.
3. You could believe you're executing queries on a single connection, inside of a transaction, when in reality Go has created several connections for you invisibly and some statements aren't part of the transaction.

{% include toc.md %}
