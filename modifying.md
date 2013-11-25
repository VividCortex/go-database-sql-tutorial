---
layout: page
permalink: /modifying/
title: Statements that Modify Data
tags: 
image:
  feature: abstract-5.jpg
share: true
---

As mentioned previously, you should only use `Query` functions to execute queries -- statements that return rows. Use `Exec`, preferably with a prepared statement, to accomplish an INSERT, UPDATE, DELETE, or other statement that doesn't return rows. The following example shows how to insert a row:

```go
stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
if err != nil {
	log.Fatal(err)
}
res, err := stmt.Exec("Dolly")
if err != nil {
	log.Fatal(err)
}
lastId, err := res.LastInsertId()
if err != nil {
	log.Fatal(err)
}
rowCnt, err := res.RowsAffected()
if err != nil {
	log.Fatal(err)
}
log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
```

Executing the statement produces a `sql.Result` that gives access to statement metadata: the last inserted ID and the number of rows affected.

What if you don't care about the result? What if you just want to execute a statement and check if there were any errors, but ignore the result? Wouldn't the following two statements do the same thing?

```go
_, err := db.Exec("DELETE FROM users")  // OK
_, err := db.Query("DELETE FROM users") // BAD
```

The answer is no. They do **not** do the same thing, and **you should never use `Query()` like this.** The `Query()` will return a `sql.Rows`, which will not be released until it's garbage collected, which can be a long time. During that time, it will continue to hold open the underlying connection, and this anti-pattern is therefore a good way to run out of resources (too many connections, for example).

{% include toc.md %}
