---
layout: article
title: Importing a Database Driver
---

To use `database/sql` you'll need the package itself as well as a driver for the specific database you want to use.

It is considered best practice to refer only to the types defined in `database/sql` as opposed to accessing drivers directly. This avoids creating code dependent on the driver allowing you to swap drivers and/or the type of database with minimal code changes on your part. It also relies on consistent Go idioms instead of purpose built ones that a particular driver author may have created to suit their own context.

In this documentation we'll use the excellent [MySQL drivers](https://github.com/go-sql-driver/mysql) from @julienschmidt and @arnehormann for examples. You can [find other drivers here](https://github.com/golang/go/wiki/SQLDrivers).

Add the following to the top of your Go source file:

<pre class="prettyprint lang-go">
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
</pre>

Notice that we're loading the driver anonymously; aliasing its package qualifier to `_` so none of its exported names are visible to our code. Under the hood the driver registers itself as being available to the `database/sql` package, but in general nothing else happens.

Now you're ready to access a database.

**Previous: [Overview of Go's database/sql Package](overview.html)**
**Next: [Accessing the Database](accessing.html)**
