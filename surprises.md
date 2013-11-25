---
layout: page
permalink: /surprises/
title: Surprises, Antipatterns and Limitations
tags: 
image:
  feature: abstract-5.jpg
share: true
---

We've documented several surprises and antipatterns throughout this tutorial, so please refer back to them if you didn't read them already:

* Opening and closing databases can cause exhaustion of resources.
* Failing to use `rows.Close()` or `stmt.Close()` can cause exhaustion of resources.
* Mixing transactional methods and transactional SQL commands can cause exhaustion of resources.
* Using `Query()` for a statement that doesn't return rows is a bad idea.
* Failing to use prepared statements can lead to a lot of extra database activity.
* Nulls cause annoying problems, which may show up only in production.
* There's a special error when there's an empty result set, which should be checked to avoid application bugs.

There are also a couple of limitations in the `database/sql` package. The interface doesn't give you all-encompassing access to what's happening under the hood. For example, you don't have much control over the pool of connections.

Another limitation, which can be a surprise, is that you can't pass big unsigned integers as parameters to statements if their high bit is set:

```go
_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64)
```

This will throw an error. Be careful if you use `uint64` values, as they may start out small and work without error, but increment over time and start throwing errors.

Another possible surprise is the effects of doing things that change connection state, such as
setting the current database with a `USE` statement. This will affect only the connection
that you run it in. Unless you are in a transaction, other statements that you think are
executed on that connection may actually run on different connections gotten from the pool, so they won't see
the effects of such changes. Additionally, after you've changed the connection, it'll return
to the pool and potentially pollute the state for some other code. This is one of the reasons
why you should never issue BEGIN or COMMIT statements as SQL commands directly, too.

{% include toc.md %}
