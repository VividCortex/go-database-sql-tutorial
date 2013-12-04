---
layout: article
title: Retrieving Result Sets
---

There are several idiomatic operations to retrieve results from the datastore.

1. Execute a query that returns rows.
1. Prepare a statement for repeated use, execute it multiple times, and destroy it.
1. Execute a statement in a once-off fashion, without preparing it for repeated use.
1. Execute a query that returns a single row. There is a shortcut for this special case.

Go's `database/sql` function names are significant. **If a function name
includes `Query`, it is designed to ask a question of the database, and will
return a set of rows**, even if it's empty. Statements that don't return rows
should not use `Query` functions; they should use `Exec()`.

Fetching Data from the Database
===============================

Let's take a look at an example of how to query the database, working with
results. We'll query the `users` table for a user whose `id` is 1, and print out
the user's `id` and `name`.  We will assign results to variables, a row at a
time, with `rows.Scan()`.

<pre class="prettyprint lang-go">
	var (
		id int
		name string
	)
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
</pre>

Here's what's happening in the above code:

1. We're using `db.Query()` to send the query to the database. We check the error, as usual.
2. We defer `rows.Close()`. This is very important.
3. We iterate over the rows with `rows.Next()`.
4. We read the columns in each row into variables with `rows.Scan()`.
5. We check for errors after we're done iterating over the rows.

A couple parts of this are easy to get wrong, and can have bad consequences.

First, as long as there's an open result set (represented by `rows`), the
underlying connection is busy and can't be used for any other query. That means
it's not available in the connection pool. If you iterate over all of the rows
with `rows.Next()`, eventually you'll read the last row, and `rows.Next()` will
encounter an internal EOF error and call `rows.Close()` for you. But if for any
reason you exit that loop -- an error, an early return, or so on -- then the
`rows` doesn't get closed, and the connection remains open. This is an easy way
to run out of resources. This is why **you should always `defer rows.Close()`**,
even if you also call it explicitly at the end of the loop, which isn't a bad
idea. `rows.Close()` is a harmless no-op if it's already closed, so you can call
it multiple times. Notice, however, that we check the error first, and only do
`rows.Close()` if there isn't an error, in order to avoid a runtime panic.

Second, you should always check for an error at the end of the `for rows.Next()`
loop. If there's an error during the loop, you need to know about it. Don't just
assume that the loop iterates until you've processed all the rows.

The error returned by `rows.Close()` is the only exception to the general rule
that it's best to capture and check for errors in all database operations. If
`rows.Close()` throws an error, it's unclear what is the right thing to do.
Logging the error message or panicing might be the only sensible thing to do,
and if that's not sensible, then perhaps you should just ignore the error.

This is pretty much the only way to do it in Go. You can't
get a row as a map, for example. That's because everything is strongly typed.
You need to create variables of the correct type and pass pointers to them, as
shown.

Preparing Queries
=================

You should, in general, always prepare queries to be used multiple times. The
result of preparing the query is a prepared statement, which can have
placeholders (a.k.a. bind values) for parameters that you'll provide when you
execute the statement.  This is much better than concatenating strings, for all
the usual reasons (avoiding SQL injection attacks, for example).

In MySQL, the parameter placeholder is `?`, and in PostgreSQL it is `$N`, where
N is a number. In Oracle placeholders begin with a colon and are named, like
`:param1`. We'll use `?` because we're using MySQL as our example.

<pre class="prettyprint lang-go">
	stmt, err := db.Prepare("select id, name from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		// ...
	}
</pre>

Under the hood, `db.Query()` actually prepares, executes, and closes a prepared
statement. That's three round-trips to the database. If you're not careful, you
can triple the number of database interactions your application makes! Some
drivers can avoid this in specific cases with an addition to `database/sql` in
Go 1.1, but not all drivers are smart enough to do that. Caveat Emptor.

Naturally prepared statements and the managment of prepared statements cost
resources. You should take care to close statements when they are not used again.

Single-Row Queries
==================

If a query returns at most one row, you can use a shortcut around some of the
lengthy boilerplate code:

<pre class="prettyprint lang-go">
	var name string
	err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
</pre>

Errors from the query are deferred until `Scan()` is called, and then are
returned from that. You can also call `QueryRow()` on a prepared statement:

<pre class="prettyprint lang-go">
	stmt, err := db.Prepare("select id, name from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var name string
	err = stmt.QueryRow(1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
</pre>

Go defines a special error constant, called `sql.ErrNoRows`, which is returned
from `QueryRow()` when the result is empty. This needs to be handled as a
special case in most circumstances. An empty result is often not considered an
error by application code, and if you don't check whether an error is this
special constant, you'll cause application-code errors you didn't expect.

One might ask why an empty result set is considered an error. There's nothing
erroneous about an empty set. The reason is that the `QueryRow()` method needs
to use this special-case in order to let the caller distinguish whether
`QueryRow()` in fact found a row; without it, `Scan()` wouldn't do anything and
you might not realize that your variable didn't get any value from the database
after all.

You should not run into this error when you're not using `QueryRow()`. If you
encounter this error elsewhere, you're doing something wrong.

**Previous: [Accessing the Database](accessing.html)**
**Next: [Modifying Data and Using Transactions](modifying.html)**
