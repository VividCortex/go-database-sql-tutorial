---
layout: page
permalink: /accessing/
title: Accessing the Database
tags: 
image:
  feature: abstract-5.jpg
share: true
---


Now that you've loaded the driver package, you're ready to create a database
object, a `sql.DB`. The first thing you should know is that **a `sql.DB` isn't a
database connection**. It also doesn't map to any particular database software's
notion of a "database" or "schema." It's an abstraction of the interface and
existence of a database, which might be a local file, accessed through a network
connection, in-memory and in-process, or what have you.

The `sql.DB` performs some important tasks for you behind the scenes:

* It opens and closes connections to the actual underlying database, via the driver.
* It manages a pool of connections as needed.

These "connections" may be file handles, sockets, network connections, or other ways to access the database. The `sql.DB` abstraction is designed to keep you from worrying about how to manage concurrent access to the underlying datastore. A connection is marked in-use when you use it to perform a task, and then returned to the available pool when it's not in use anymore. One consequence of this is that **if you fail to release connections back to the pool, you can cause `db.SQL` to open a lot of connections**, potentially running out of resources (too many connections, too many open file handles, lack of available network ports, etc). We'll discuss more about this later.

To create a `sql.DB`, you use `sql.Open()`. This returns a `*sql.DB`:

```go
func main() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
```

In the example shown, we're illustrating several things:

1. The first argument to `sql.Open` is the driver name. This is the string that the driver used to register itself with `database/sql`, and is conventionally the same as the package name to avoid confusion.
2. The second argument is a driver-specific syntax that tells the driver how to access the underlying datastore. In this example, we're connecting to the "hello" database inside a local MySQL server instance.
3. You should (almost) always check and handle errors returned from all `database/sql` operations.
4. It is idiomatic to `defer db.Close()` if the `sql.DB` should not have a lifetime beyond the scope of the function.

Perhaps counter-intuitively, `sql.Open()` **does not establish any connections to the database**, nor does it validate driver connection parameters. Instead, it simply prepares the database abstraction for later use. The first actual connection to the underlying datastore will be established lazily, when it's needed for the first time. If you want to check right away that the database is available and accessible (for example, check that you can establish a network connection and log in), use `db.Ping()` to do that, and remember to check for errors:

```go
	err = db.Ping()
	if err != nil {
		// do something here
	}
```

Although it's idiomatic to `Close()` the database when you're finished with it, **the `sql.DB` object is designed to be long-lived.** Don't `Open()` and `Close()` databases frequently. Instead, create **one** `sql.DB` object for each distinct datastore you need to access, and keep it until the program is done accessing that datastore. Pass it around as needed, or make it available somehow globally, but keep it open. And don't `Open()` and `Close()` from a short-lived function. Instead, pass the `sql.DB` into that short-lived function as an argument.

If you don't treat the `sql.DB` as a long-lived object, you could experience problems such as poor reuse and sharing of connections, running out of available network resources, sporadic failures due to a lot of TCP connections remaining in TIME_WAIT status, and so on.

{% include toc.md %}
