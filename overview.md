---
layout: page
permalink: /overview/
title: Overview
tags: 
image:
  feature: abstract-5.jpg
share: false
---

To access databases in Go, you use a `sql.DB`. You use this type to create
statements and transactions, execute queries, and fetch results.

The first thing you should know is that **a `sql.DB` isn't a database
connection**. It also doesn't map to any particular database software's notion
of a "database" or "schema." It's an abstraction of the interface and existence
of a database, which might be as varied as a local file, accessed through a network
connection, or in-memory and in-process.

The `sql.DB` performs some important tasks for you behind the scenes:

* It opens and closes connections to the actual underlying database, via the driver.
* It manages a pool of connections as needed, which may be a variety of things as mentioned.

The `sql.DB` abstraction is designed to keep you from worrying about how to
manage concurrent access to the underlying datastore.  A connection is marked
in-use when you use it to perform a task, and then returned to the available
pool when it's not in use anymore. One consequence of this is that **if you fail
to release connections back to the pool, you can cause `db.SQL` to open a lot of
connections**, potentially running out of resources (too many connections, too
many open file handles, lack of available network ports, etc). We'll discuss
more about this later.

After creating a `sql.DB`, you can use it to query the database that it
represents, as well as creating statements and transactions.

{% include toc.md %}
