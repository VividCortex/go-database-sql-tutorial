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
* There's no limit on the number of connections. If you try to do a lot of things at once, you can create an arbitrary number of connections. This can cause the database to return an error such as "too many connections."
* In Go 1.1, you can use `db.SetMaxIdleConns(N)` to limit the number of *idle* connections in the pool. This doesn't limit the pool size, though. That will probably come in Go 1.3 or later.
* Connections are recycled rather fast. Setting the number of idle connections with `db.SetMaxIdleConns(N)` can reduce this churn, and help keep connections around for reuse.

In a future version of Go, the `database/sql` package will have the ability to
[limit open connections](https://code.google.com/p/go/issues/detail?id=4805). In
the meantime, it's still rather easy to run out of database connections by
exceeding the limit set in the database itself (e.g. `max_connections` in
MySQL). To avoid this, you'll need to limit your application's concurrent use of
the `sql.DB` object. The idiomatic way to do that is described in the [channels
section of the Effective Go](http://golang.org/doc/effective_go.html#channels)
document.

{% include toc.md %}
