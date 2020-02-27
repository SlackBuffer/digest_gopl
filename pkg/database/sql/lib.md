https://golang.org/pkg/database/sql/  
https://github.com/golang/go/wiki/SQLInterface  
https://godoc.org/github.com/lib/pq

- [ ] https://play.golang.org/p/wLD2ykbJUMr
- [ ] https://golang.org/src/database/sql/example_test.go

`DB` is a database handle representing a pool of zero or more underlying connections. It's **safe for concurrent use** by multiple goroutines.

The `sql` package creates and frees connections automatically; it also maintains a free pool of idle connections. If the database has a concept of per-connection state, such state can be reliably observed within a transaction (`Tx`) or connection (`Conn`). Once `DB.Begin` is called, the returned `Tx` is bound to a single connection. Once `Commit` or `Rollback` is called on the transaction, that transaction's connection is returned to DB's idle connection pool. The pool size can be controlled with `SetMaxIdleConns`.

`Open` may just validate its arguments without creating a connection to the database. To verify that the data source name is valid, call `Ping`.  
The returned `DB` is safe for **concurrent** use by multiple goroutines and maintains its own pool of idle connections. Thus, the `Open` function should be called just **once**. 

It is rare to `Close` a `DB`, as the `DB` handle is meant to be long-lived and shared between many goroutines.

A `*DB` is a **pool** of connections. Call `Conn` to reserve a connection for exclusive use.

`Ping` and `PingContext` may be used to determine if communication with the database server is still possible. When used in a command line application `Ping` may be used to establish that further queries are possible; that the provided DSN is valid. When used in long running service `Ping` may be part of the **health checking** system.

If the database is being written to, ensure to check for `Close` errors that may be returned from the driver. The query may encounter an auto-commit error and be forced to rollback changes.

In normal use, create one `Stmt` when your process starts. Then reuse it each time you need to issue the query.

```go
defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
```

