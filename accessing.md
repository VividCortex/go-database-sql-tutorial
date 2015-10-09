---
layout: article
title: Accessing the Database
---

Now that you've loaded the driver package, you're ready to create a database
object, a `sql.DB`.

To create a `sql.DB`, you use `sql.Open()`. This returns a `*sql.DB`:

<pre class="prettyprint lang-go">
func main() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
</pre>

In the example shown, we're illustrating several things:

1. The first argument to `sql.Open` is the driver name. This is the string that the driver used to register itself with `database/sql`, and is conventionally the same as the package name to avoid confusion. For example, it's `mysql` for [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql). Some drivers do not follow the convention and use the database name, e.g. `sqlite3` for [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) and `postgres` for [github.com/lib/pq](https://github.com/lib/pq).
2. The second argument is a driver-specific syntax that tells the driver how to access the underlying datastore. In this example, we're connecting to the "hello" database inside a local MySQL server instance.
3. You should (almost) always check and handle errors returned from all `database/sql` operations.  There are a few special cases that we'll discuss later where it doesn't make sense to do this.
4. It is idiomatic to `defer db.Close()` if the `sql.DB` should not have a lifetime beyond the scope of the function.

Perhaps counter-intuitively, `sql.Open()` **does not establish any connections
to the database**, nor does it validate driver connection parameters. Instead,
it simply prepares the database abstraction for later use. The first actual
connection to the underlying datastore will be established lazily, when it's
needed for the first time. If you want to check right away that the database is
available and accessible (for example, check that you can establish a network
connection and log in), use `db.Ping()` to do that, and remember to check for
errors:

<pre class="prettyprint lang-go">
err = db.Ping()
if err != nil {
	// do something here
}
</pre>

Although it's idiomatic to `Close()` the database when you're finished with it,
**the `sql.DB` object is designed to be long-lived.** Don't `Open()` and
`Close()` databases frequently. Instead, create **one** `sql.DB` object for each
distinct datastore you need to access, and keep it until the program is done
accessing that datastore. Pass it around as needed, or make it available somehow
globally, but keep it open. And don't `Open()` and `Close()` from a short-lived
function. Instead, pass the `sql.DB` into that short-lived function as an
argument.

If you don't treat the `sql.DB` as a long-lived object, you could experience
problems such as poor reuse and sharing of connections, running out of available
network resources, or sporadic failures due to a lot of TCP connections
remaining in `TIME_WAIT` status. Such problems are signs that you're not using
`database/sql` as it was designed.

Now it's time to use your `sql.DB` object.

**Previous: [Importing a Database Driver](importing.html)**
**Next: [Retrieving Result Sets](retrieving.html)**
