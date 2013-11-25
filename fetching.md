---
layout: page
permalink: /fetching/
title: Fetching Data from teh Database
tags: 
image:
  feature: abstract-5.jpg
share: true
---

Let's take a look at an example of how to query the database, working with results. We'll query the `users` table for a user whose `id` is 1, and print out the user's `id` and `name`:

```go
	var (
		id int
		name string
	)
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
```

Here's what's happening in the above code:

1. We're using `db.Query()` to send the query to the database. We check the error, as usual.
2. We defer `rows.Close()`. This is very important; more on that in a moment.
3. We iterate over the rows with `rows.Next()`.
4. We read the columns in each row into variables with `rows.Scan()`.
5. We check for errors after we're done iterating over the rows.

A couple parts of this are easy to get wrong, and can have bad consequences.

First, as long as there's an open result set (represented by `rows`), the underlying connection is busy and can't be used for any other query. That means it's not available in the connection pool. If you iterate over all of the rows with `rows.Next()`, eventually you'll read the last row, and `rows.Next()` will encounter an internal EOF error and call `rows.Close()` for you. But if for any reason you exit that loop -- an error, an early return, or so on -- then the `rows` doesn't get closed, and the connection remains open. This is an easy way to run out of resources. This is why **you should always `defer rows.Close()`**, even if you also call it explicitly at the end of the loop, which isn't a bad idea. `rows.Close()` is a harmless no-op if it's already closed, so you can call it multiple times. Notice, however, that we check the error first, and only do `rows.Close()` if there isn't an error, in order to avoid a runtime panic.

Second, you should always check for an error at the end of the `for rows.Next()` loop. If there's an error during the loop, you need to know about it. Don't just assume that the loop iterates until you've processed all the rows.

The error returned by `rows.Close()` is the only exception to the general rule that it's best to capture and check for errors in all database operations. If `rows.Close()` throws an error, it's unclear what is the right thing to do. Logging the error message or panicing might be the only sensible thing to do, and if that's not sensible, then perhaps you should just ignore the error.

{% include toc.md %}
