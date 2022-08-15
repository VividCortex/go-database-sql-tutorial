---
layout: article
title: Overview
---

To access databases in Go, you use a `sql.DB`. You use this type to create
statements and transactions, execute queries, and fetch results.

The first thing you should know is that **a `sql.DB` isn't a database
connection**. It also doesn't map to any particular database software's notion
of a "database" or "schema." It's an abstraction of the interface and existence
of a database, which might be as varied as a local file, accessed through a network
connection, or in-memory and in-process.

`sql.DB` differs from equivalent entity in C#, Java or PHP. It creates a handle with empty slice of connections.
Subsequent DB operation (Query, QueryRow, Execute, Ping) creates real conection and pushes it into the slice. The connection 
is marked as 'free' after finishing an operation and subsequent one can reuse it instead of creating new connection.
Query is finished after closing 'Rows'. QueryRow, Execute or Ping closes connection immediately. DB object allows some
optimization of the pool calling functions SetMaxOpenConns(), SetMaxIdleConns() and SetConnMaxLifetime(). Details of the optimization are described 
on site 'https://www.alexedwards.net/blog/configuring-sqldb'.

The `sql.DB` abstraction is designed to keep you from worrying about how to
manage concurrent access to the underlying datastore.  A connection is marked
in-use when you use it to perform a task, and then returned to the available
pool when it's not in use anymore. One consequence of this is that **if you fail
to release connections back to the pool, you can cause `sql.DB` to open a lot of
connections**, potentially running out of resources (too many connections, too
many open file handles, lack of available network ports, etc). We'll discuss
more about this later.

After creating a `sql.DB`, you can use it to query the database that it
represents, as well as creating statements and transactions.

**Next: [Importing a Database Driver](importing.html)**
