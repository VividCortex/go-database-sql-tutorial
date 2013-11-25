---
layout: page
permalink: /operations/
title: Common Database Operations
tags: 
image:
  feature: abstract-5.jpg
share: true
---

Now that you've loaded the driver and opened the `sql.DB`, it's time to use it. There are several idiomatic operations against the datastore:

1. Execute a query that returns rows.
1. Execute a query that returns a single row. There is a shortcut for this special case.
1. Prepare a statement for repeated use, execute it multiple times, and destroy it.
1. Execute a statement in a once-off fashion, without preparing it for repeated use.
1. Modify data and check for the results.
1. Perform transaction-related operations; not discussed at this time.

You should almost always capture and check errors from all functions that return them. There are a few special cases that we'll discuss later where it doesn't make sense to do this.

Go's `database/sql` function names are significant. **If a function name includes `Query`, it is designed to ask a question of the database, and should return a set of rows**, even if it's empty. Statements that don't return rows should not use `Query` functions, for reasons we'll also discuss later.

{% include toc.md %}
