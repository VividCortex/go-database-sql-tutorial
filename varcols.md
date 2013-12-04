---
layout: article
title: Working with Unknown Columns
---

The `Scan()` function requires you to pass exactly the right number of
destination variables. What if you don't know what the query will return?

If you don't know how many columns the query will return, you can use
`Columns()` to find a list of column names. You can examine the length of this
list to see how many columns there are, and you can pass a slice into `Scan()`
with the correct number of values. For example, some forks of MySQL return
different columns for the `SHOW PROCESSLIST` command, so you have to be prepared
for that or you'll cause an error. Here's one way to do it; there are others:

<pre class="prettyprint lang-go">
cols, err := rows.Columns()
if err != nil {
	// handle the error
} else {
	dest := []interface{}{ // Standard MySQL columns
		new(uint64), // id
		new(string), // host
		new(string), // user
		new(string), // db
		new(string), // command
		new(uint32), // time
		new(string), // state
		new(string), // info
	}
	if len(cols) == 11 {
		// Percona Server
	} else if len(cols) &gt; 8 {
		// Handle this case
	}
	err = rows.Scan(dest...)
	// Work with the values in dest
}
</pre>

If you don't know the columns or their types, you should use `sql.RawBytes`.

<pre class="prettyprint lang-go">
cols, err := rows.Columns() // Remember to check err afterwards
vals := make([]interface{}, len(cols))
for i, _ := range cols {
	vals[i] = new(sql.RawBytes)
}
for rows.Next() {
	err = rows.Scan(vals...)
	// Now you can check each element of vals for nil-ness,
	// and you can use type introspection and type assertions
	// to fetch the column into a typed variable.
}
</pre>

**Previous: [Working with NULLs](nulls.html)**
**Next: [The Connection Pool](connection-pool.html)**
