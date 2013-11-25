---
layout: page
permalink: /preparing/
title: Preparing Queries
tags: 
image:
  feature: abstract-5.jpg
share: true
---

You should, in general, always prepare queries to be used multiple times. The result of preparing the query is a prepared statement, which can have `?` placeholders for parameters that you'll provide when you execute the statement. This is much better than concatenating strings, for all the usual reasons (avoiding SQL injection attacks, for example).

```go
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
```

Under the hood, `db.Query()` actually prepares, executes, and closes a prepared statement. That's three round-trips to the database. If you're not careful, you can triple the number of database interactions your application makes! Some drivers can avoid this in specific cases with an addition to `database/sql` in Go 1.1, but not all drivers are smart enough to do that. Caveat Emptor.

Statements are like results: they claim a connection and should be closed. It's idiomatic to `defer stmt.Close()` if the prepared statement `stmt` should not have a lifetime beyond the scope of the function. If you don't, it'll reserve a connection from the pool.

{% include toc.md %}
