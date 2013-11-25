---
layout: page
permalink: /single-row-queries/
title: Single-Row Queries
tags: 
image:
  feature: abstract-5.jpg
share: true
---

If a query returns at most one row, you can use a shortcut around some of the lengthy boilerplate code:

```go
var name string
err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
if err != nil {
	log.Fatal(err)
}
fmt.Println(name)
```

Errors from the query are deferred until `Scan()` is called, and then are returned from that. You can also call `QueryRow()` on a prepared statement:

```go
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
```

{% include toc.md %}
