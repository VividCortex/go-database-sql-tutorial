---
layout: page
permalink: /connection-pool/
title: The Connection Pool
tags:
image:
  feature: abstract-5.jpg
share: false
---

There is a basic connection pool in the `database/sql` package. There isn't a
lot of ability to control or inspect it, but here are some things you might find
useful to know:

* Connections are created when needed and there isn't a free connection in the pool.
* By default, there's no limit on the number of connections. If you try to do a lot of things at once, you can create an arbitrary number of connections. This can cause the database to return an error such as "too many connections."
* In Go 1.1 or newer, you can use `db.SetMaxIdleConns(N)` to limit the number of *idle* connections in the pool. This doesn't limit the pool size, though.
* In Go 1.2 or newer, you can use `db.SetMaxOpenConns(N)` to set the maximum of *total* open connections in the pool.
* Connections are recycled rather fast. Setting a high number of idle connections with `db.SetMaxIdleConns(N)` can reduce this churn, and help keep connections around for reuse.

{% include toc.md %}
