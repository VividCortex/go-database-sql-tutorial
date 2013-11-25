---
layout: page
permalink: /overview/
title: Overview
tags: 
image:
  feature: abstract-5.jpg
share: true
---

The idiomatic way to use a SQL, or SQL-like, database in Go is through the
`database/sql` package. It provides a lightweight interface to a row-oriented
database. This documentation is a reference for the most common aspects of how
to use it.

The first thing to do is import the `database/sql` package, and a driver
package. You generally shouldn't use the driver package directly, although some
drivers encourage you to do so. (In our opinion, it's usually a bad idea.)
Instead, your code should only refer to `database/sql`. This helps avoid making
your code dependent on the driver, so that you can change the underlying driver
(and thus the database you're accessing) without changing your code. It also
forces you to use the Go idioms instead of ad-hoc idioms that a particular
driver author may have provided.

Keep in mind that only the `database/sql` API provides this abstraction.
Specific databases and drivers can differ in behavior and/or syntax.  One
example is the syntax for placeholder parameters in prepared statements. For
example, comparing MySQL, PostgreSQL, and Oracle:

	MySQL               PostgreSQL            Oracle
	=====               ==========            ======
	WHERE col = ?       WHERE col = $1        WHERE col = :col
	VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)

In this documentation, we'll use the excellent
[MySQL drivers](https://github.com/go-sql-driver/mysql) from @arnehormann and @julienschmidt for examples.

Add the following to the top of your Go source file:

	import (
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
	)

Notice that we're loading the driver anonymously, aliasing its package qualifier
to `_` so none of its exported names are visible to our code. Under the hood,
the driver registers itself as being available to the `database/sql` package,
but in general nothing else happens.

Now you're ready to access a database.

{% include toc.md %}
