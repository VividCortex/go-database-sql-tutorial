---
layout: page
permalink: /assigning/
title: Assigning Results to Variables
tags: 
image:
  feature: abstract-5.jpg
share: true
---

In the previous section you already saw the idiom for assigning results to variables, a row at a time, with `rows.Scan()`. This is pretty much the only way to do it in Go. You can't get a row as a map, for example. That's because everything is strongly typed. You need to create variables of the correct type and pass pointers to them, as shown.

There are two special cases: nullable columns, and variable numbers of columns, that are a little harder to handle.

Nullable columns are annoying and lead to a lot of ugly code. If you can, avoid them. If not, then you'll need to use special types from the `database/sql` package to handle them. There are types for nullable booleans, strings, integers, and floats. Here's how you use them:

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

Limitations of the nullable types, and reasons to avoid nullable columns in case you need more convincing:

1. There's no `sql.NullUint64` or `sql.NullYourFavoriteType`.
1. Nullability can be tricky, and not future-proof. If you think something won't be null, but you're wrong, your program will crash, perhaps rarely enough that you won't catch errors before you ship them.
1. One of the nice things about Go is having a useful default zero-value for every variable. This isn't the way nullable things work.

The other special case is assigning a variable number of columns into variables. The `rows.Scan()` function accepts a variable number of `interface{}`, and you have to pass the correct number of arguments. If you don't know the columns or their types, you should use `sql.RawBytes`:

	cols, err := rows.Columns()				// Get the column names; remember to check err
	vals := make([]sql.RawBytes, len(cols)) // Allocate enough values
	ints := make([]interface{}, len(cols)) 	// Make a slice of []interface{}
	for i := range ints {
		ints[i] = &vals[i] // Copy references into the slice
	}
	for rows.Next() {
		err := rows.Scan(vals...)
		// Now you can check each element of vals for nil-ness,
		// and you can use type introspection and type assertions
		// to fetch the column into a typed variable.
	}

If you know the possible sets of columns and their types, it can be a little less annoying, though still not great. In that case, you simply need to examine `rows.Columns()`, which returns an array of column names.

{% include toc.md %}
