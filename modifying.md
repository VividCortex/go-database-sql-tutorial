---
layout: page
permalink: /modifying/
title: Modifying Data and Using Transactions
tags: 
image:
  feature: abstract-5.jpg
share: false
---

Now we're ready to see how to modify data and work with transactions. The
distinction might seem artificial if you're used to programming languages that
use a "statement" object for fetching rows as well as updating data, but in Go,
there's an important reason for the difference.

Statements that Modify Data
===========================

Use `Exec()`, preferably with a prepared statement, to accomplish an `INSERT`,
`UPDATE`, `DELETE`, or other statement that doesn't return rows. The following
example shows how to insert a row and inspect metadata about the operation:

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

Executing the statement produces a `sql.Result` that gives access to statement
metadata: the last inserted ID and the number of rows affected.

What if you don't care about the result? What if you just want to execute a
statement and check if there were any errors, but ignore the result? Wouldn't
the following two statements do the same thing?

	_, err := db.Exec("DELETE FROM users")  // OK
	_, err := db.Query("DELETE FROM users") // BAD

The answer is no. They do **not** do the same thing, and **you should never use
`Query()` like this.** The `Query()` will return a `sql.Rows`, which will not be
released until it's garbage collected, which can be a long time. During that
time, it will continue to hold open the underlying connection, and this
anti-pattern is therefore a good way to run out of resources (too many
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
that transaction, and can't be used outside of it. Likewise, prepared statements
created only on a database handle can't be used within a transaction.

You should not mingle the use of transaction-related functions such as `Begin()`
and `Commit()` with SQL statements such as `BEGIN` and `COMMIT` in your SQL
code. Bad things might result:

* The `Tx` objects could remain open, reserving a connection from the pool and not returning it.
* The state of the database could get out of sync with the state of the Go variables representing it.
* You could believe you're executing queries on a single connection, inside of a transaction, when in reality Go has created several connections for you invisibly and some statements aren't part of the transaction.

{% include toc.md %}
