---
layout: page
permalink: /nulls/
title: Working with NULLs
tags: 
image:
  feature: abstract-5.jpg
share: false
---

Nullable columns are annoying and lead to a lot of ugly code. If you can, avoid
them. If not, then you'll need to use special types from the `database/sql`
package to handle them, or define your own.

There are types for nullable booleans, strings, integers, and floats. Here's how
you use them:

```go
for rows.Next() {
	var s sql.NullString
	err := rows.Scan(&s)
	// check err
	if s.Valid {
	   // use s.String
	} else {
	   // NULL value
	}
}
```

Limitations of the nullable types, and reasons to avoid nullable columns in case
you need more convincing:

1. There's no `sql.NullUint64` or `sql.NullYourFavoriteType`. You'd need to
	define your own for this.
1. Nullability can be tricky, and not future-proof. If you think something won't
	be null, but you're wrong, your program will crash, perhaps rarely enough
	that you won't catch errors before you ship them.
1. One of the nice things about Go is having a useful default zero-value for
	every variable. This isn't the way nullable things work.

If you need to define your own types to handle NULLs, you can copy the design of
`sql.NullString` to achieve that.

{% include toc.md %}
