---
layout: article
title: Importing a Database Driver
---

To use `database/sql` you'll need the package itself, as well as a driver for
the specific database you want to use.

You generally shouldn't use driver packages directly, although some drivers
encourage you to do so. (In our opinion, it's usually a bad idea.) Instead, your
code should only refer to types defined in `database/sql`, if possible. This
helps avoid making your code dependent on the driver, so that you can change the
underlying driver (and thus the database you're accessing) with minimal code
changes. It also forces you to use the Go idioms instead of ad-hoc idioms that a
particular driver author may have provided.

In this documentation, we'll use the excellent [MySQL
drivers](https://github.com/go-sql-driver/mysql) from @julienschmidt and @arnehormann
 for examples.

Add the following to the top of your Go source file:

<pre class="prettyprint lang-go">
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
</pre>

Notice that we're loading the driver anonymously, aliasing its package qualifier
to `_` so none of its exported names are visible to our code. Under the hood,
the driver registers itself as being available to the `database/sql` package,
but in general nothing else happens with the exception that the init function is run.

Now you're ready to access a database.

**Previous: [Overview of Go's database/sql Package](overview.html)**
**Next: [Accessing the Database](accessing.html)**
