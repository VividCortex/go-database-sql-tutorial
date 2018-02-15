---
layout: article
title: Working with NULLs
---

Nullable columns are annoying and lead to a lot of ugly code. If you can, avoid
them. If not, then you'll need to use special types from the `database/sql`
package to handle them, or define your own.

There are types for nullable booleans, strings, integers, and floats. Here's how
you use them:

<pre class="prettyprint lang-go">
for rows.Next() {
	var s sql.NullString
	err := rows.Scan(&amp;s)
	// check err
	if s.Valid {
	   // use s.String
	} else {
	   // NULL value
	}
}
</pre>

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

If you can't avoid having NULL values in your database, there is another work around that most database systems support, namely `COALESCE()`. Something like the following might be something that you can use, without introducing a myriad of `sql.Null*` types.
<pre class="prettyprint lang-go">
rows, err := db.Query(`
	SELECT
		name,
		COALESCE(other_field, '') as otherField
	WHERE id = ?
`, 42)

for rows.Next() {
	err := rows.Scan(&name, &otherField)
	// ..
	// If `other_field` was NULL, `otherField` is now an empty string. This works with other data types as well.
}
</pre>


**Previous: [Handling Errors](errors.html)**
**Next: [Working with Unknown Columns](varcols.html)**
