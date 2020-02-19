http://go-database-sql.org/overview.html  
https://godoc.org/github.com/lib/pq
https://www.postgresqltutorial.com/postgresql-serial/  
# Setup

```
docker run --name learn-sql --rm -p 1234:5432 -e POSTGRES_DB=testdb -e POSTGRES_USER=hofungkoeng -e POSTGRES_PASSWORD=pswd postgres

docker exec -it learn-sql psql -U hofungkoeng -d testdb
```

```sql
CREATE TABLE Product
    (product_id      CHAR(4)      NOT NULL,
     product_name    VARCHAR(100) NOT NULL,
     product_type    VARCHAR(32)  NOT NULL,
     sale_price      INTEGER,
     purchase_price  INTEGER,
     regist_date     DATE,
     id SERIAL,                 -- 放末尾以免改变以下 transaction 的变量对应关系
     PRIMARY KEY (product_id));

BEGIN TRANSACTION;
INSERT INTO Product VALUES ('0001', 'T恤衫', '衣服', 1000, 500, '2009-09-20');
INSERT INTO Product VALUES ('0002', '打孔器', '办公用品', 500, 320, '2009-09-11');
INSERT INTO Product VALUES ('0003', '运动T恤', '衣服', 4000, 2800, NULL);
INSERT INTO Product VALUES ('0004', '菜刀', '厨房用具', 3000, 2800, '2009-09-20');
INSERT INTO Product VALUES ('0005', '高压锅', '厨房用具', 6800, 5000, '2009-01-15');
INSERT INTO Product VALUES ('0006', '叉子', '厨房用具', 500, NULL, '2009-09-20');
INSERT INTO Product VALUES ('0007', '擦菜板', '厨房用具', 880, 790, '2008-04-28');
INSERT INTO Product VALUES ('0008', '圆珠笔', '办公用品', 100, NULL,'2009-11-11');
COMMIT;
```

# Overview
To access databases in Go, you use a `sql.DB` to create statements and transactions, execute queries, and fetch results.

A `sql.DB` **isn’t** a database connection. It also doesn’t map to any particular database software’s notion of a “database” or “schema.” It’s an **abstraction** of the interface and existence of a database, which might be as varied as a local file, accessed through a network connection, or in-memory and in-process.

`sql.DB` performs some important tasks behind the scenes:
- It opens and closes connections to the actual underlying database, via the **driver**.
- It manages a pool of connections as needed.

The `sql.DB` abstraction is designed to keep you from worrying about how to manage concurrent access to the underlying datastore. A connection is marked in-use when you use it to perform a task, and then returned to the available pool when it’s not in use anymore. One consequence of this is that if you fail to release connections back to the pool, you can cause `sql.DB` to open a lot of connections, potentially running out of resources (too many connections, too many open file handles, lack of available network ports, etc).
# Importing a Database Driver
We’re loading the driver anonymously, aliasing its package qualifier to `_` so none of its exported names are visible to our code. Under the hood, the driver registers itself as being available to the `database/sql` package, but in general nothing else happens with the exception that the `init` function is run.

```go
import (
	"database/sql"
	_ "github.com/lib/pq"
)
```

# Accessing the Database
Postgres connection strings
- https://godoc.org/github.com/lib/pq
- https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-CONNSTRING

```go
// DSN: database (data) source name
// connStr := "dbname=testdb user=hofungkoeng password=pswd port=1234 sslmode=disable"
connStr := "postgres://hofungkoeng:pswd@localhost:1234/testdb?sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil {
    log.Fatal("sql open error:", err)
}
defer db.Close()
if err := db.Ping(); err != nil {
    log.Fatal("db ping error:", err)
}
```

1. The first argument is the driver name. It's the string that the driver used to register itself with `database/sql`.  
It's conventionally the same as the package name to avoid confusion. For example, it’s `mysql` for `github.com/go-sql-driver/mysql`. Some drivers do not follow the convention and use the database name, e.g. `sqlite3` for `github.com/mattn/go-sqlite3` and `postgres` for `github.com/lib/pq`.
1. You should (almost) always check and handle errors returned from all `database/sql` operations.
2. It's idiomatic to defer `db.Close()` if the `sql.DB` should not have a lifetime beyond the scope of the function.
3. **\*** `sql.Open()` ***does not*** establish any connections to the database, nor does it validate driver connection parameters. Instead, it simply prepares the database abstraction for later use. The first actual connection to the underlying datastore will be established **lazily**, when it’s needed for the first time.  
If you want to check right away that the database is available and accessible (for example, check that you can establish a network connection and log in), use `db.Ping()` to do that.

Although it’s idiomatic to `Close()` the database when you’re finished with it, the `sql.DB` object is designed to be **long-lived**. Don’t `Open()` and `Close()` databases frequently. Instead, create one `sql.DB` object for each distinct datastore you need to access, and keep it until the program is done accessing that datastore. **Pass** it around as needed, or make it available somehow **globally**, but **keep it open**. And don’t `Open()` and `Close()` from a short-lived function. Instead, pass the `sql.DB` into that short-lived function as an argument.

If you don’t treat the `sql.DB` as a long-lived object, you could experience problems such as poor reuse and sharing of connections, running out of available network resources, or sporadic failures due to a lot of TCP connections remaining in `TIME_WAIT` status. Such problems are signs that you’re not using `database/sql` as it was designed.
# Retrieving Result Sets
Go’s `database/sql` function names are significant. If a function name includes `Query`, it is designed to ask a question of the database, and will return a set of rows, even if it’s empty. Statements that don’t return rows should not use `Query` functions; they should use `Exec()`.

```go
rows, err := db.Query("select product_id, product_name from product where product_id = $1", "0001")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()
for rows.Next() {
    var i item
    err := rows.Scan(&i.productID, &i.productName)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(i)
}
if err = rows.Err(); err != nil {
    log.Fatal(err)
}
```

This is pretty much the only way to do it in Go. You can’t get a row as a map, for example. That’s because everything is strongly typed. You need to create variables of the correct type and pass pointers to them.

Some notes:
1. Always check for an error at the end of the for `rows.Next()` loop. If there’s an error during the loop, you need to know about it. Don’t just assume that the loop iterates until you’ve processed all the rows.
2. As long as there’s an open result set (represented by `rows`), the underlying connection is busy and can’t be used for any other query. That means it’s not available in the connection pool. If you iterate over all of the rows with `rows.Next()`, eventually you’ll read the last row, and `rows.Next()` will encounter an internal `EOF` error and call `rows.Close()` for you. But if for some reason you exit that loop – an early return, or so on – then the rows doesn’t get closed, and the connection remains open. (It is auto-closed if `rows.Next()` returns false due to an error, though). This is an easy way to run out of resources.  
`rows.Close()` is a **harmless no-op** if it’s already closed, so you can call it multiple times. Notice, however, that we **check the error first**, and only call `rows.Close()` if there isn’t an error, in order to avoid a runtime panic.
3. Don’t `defer` within a loop. A deferred statement doesn’t get executed until the function exits, so a long-running function shouldn’t use it. If you do, you will slowly accumulate memory. If you are repeatedly querying and consuming result sets within a loop, you should explicitly call `rows.Close()` when you’re done with each result, and **not** use `defer`.

When you iterate over rows and scan them into destination variables, Go performs data type conversions work for you, behind the scenes. It is based on the type of the destination variable. This can clean up your code and help avoid repetitive work.
- For example, suppose you select some rows from a table that is defined with string columns, such as `VARCHAR(45)` or similar. You happen to know, however, that the table always contains numbers. If you pass a pointer to a string, Go will copy the bytes into the string. Now you can use `strconv.ParseInt()` or similar to convert the value to a number. You’ll have to check for errors in the SQL operations, as well as errors parsing the integer. This is messy and tedious.
- Or, you can just pass `Scan()` a pointer to an integer. Go will detect that and call `strconv.ParseInt()` for you. If there’s an error in conversion, the call to `Scan()` will return it. Your code is neater and smaller now. This is the recommended way to use `database/sql`.

```go
// statement
stmt, err := db.Prepare("select product_id, product_name from product where product_id = $1")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()
rows, err := stmt.Query("0001")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()
for rows.Next() {
    var i item
    err := rows.Scan(&i.productID, &i.productName)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(i)
}
if err = rows.Err(); err != nil {
    log.Fatal(err)
}
```

Under the hood, `db.Query()` actually prepares, executes, and closes a prepared statement. That’s **three round-trips** to the database. If you’re not careful, you can triple (because of **re-preparing**) the number of database interactions your application makes. Some drivers can avoid this in specific cases, but not all drivers do.

You should, in general, always prepare queries to be used multiple times. The result of preparing the query is a *prepared statement*, which can have **placeholders** for parameters that you’ll provide when you execute the statement. This is much better than concatenating strings, for all the usual reasons (avoiding SQL injection attacks, for example).  
In MySQL, the parameter placeholder is `?`, and in PostgreSQL it is `$N`, where `N` is a number. SQLite accepts either of these. In Oracle placeholders begin with a colon and are named, like `:param1`.

```
MySQL               PostgreSQL            Oracle
=====               ==========            ======
WHERE col = ?       WHERE col = $1        WHERE col = :col
VALUES(?, ?, ?)     VALUES($1, $2, $3)    VALUES(:val1, :val2, :val3)
```

```go
var i item
if err := db.QueryRow("select product_id, product_name from product where product_id = $1", "0001").Scan(&i.productID, &i.productName); err != nil {
    log.Fatal(err)
}
fmt.Println(i)


var i item
stmt, err := db.Prepare("select product_id, product_name from product where product_id = $1")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()
err = stmt.QueryRow("0001").Scan(&i.productID, &i.productName)
if err != nil {
    log.Fatal(err)
}
fmt.Println(i)
```

If a query returns at most one row, you can use a shortcut. Errors from the query are deferred until `Scan()` is called, and then are returned from that.  
You can also call `QueryRow()` on a prepared statement.
# Modifying Data and Using Transactions
Use `Exec()`, preferably with a prepared statement, to accomplish an `INSERT`, `UPDATE`, `DELETE`, or another statement that doesn’t return rows.

```go
stmt, err := db.Prepare("INSERT INTO product (product_id, product_name, product_type, sale_price, purchase_price, regist_date) VALUES ($1, $2, $3, $4, $5, $6)")
if err != nil {
    log.Fatal(err)
}
res, err := stmt.Exec("0016", "T恤衫", "衣服", 1000, 500, "2009-09-20")
if err != nil {
    log.Fatal(err)
}
// RowsAffected returns the number of rows affected by an update, insert, or delete.
rowCnt, err := res.RowsAffected()
if err != nil {
    log.Fatal(err)
}
log.Printf("affected row = %d\n", rowCnt)

// pq does not support the LastInsertId() method of the Result type in database/sql. To return the identifier of an INSERT (or UPDATE or DELETE), use the Postgres RETURNING clause with a standard Query or QueryRow call.
// id 需在 CREATE TABLE 时定义过
var lastID int
err = db.QueryRow(`INSERT INTO product (product_id, product_name, product_type) VALUES ('0020', 'T恤衫', '衣服') RETURNING id`).Scan(&lastID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("last ID: %d\n", lastID)

// LastInsertId returns the integer generated by the database in response to a command.
// lastID, err := res.LastInsertId()
// if err != nil {
// 	log.Fatal(err)
// }
```

Executing the statement produces a `sql.Result` that gives access to statement metadata.

```go
_, err := db.Exec("DELETE FROM product")  // OK
_, err := db.Query("DELETE FROM product") // BAD
```

The `Query()` will return a `sql.Rows`, which ***reserves*** a database connection **until** the `sql.Rows` is closed. The garbage collector will eventually close the underlying `net.Conn` for you, but this might take a long time.  
Moreover the `database/sql` package keeps tracking the connection in its pool, hoping that you release it at some point, so that the connection can be used again. This anti-pattern is therefore a good way to run out of resources (too many connections, for example).

In Go, a *transaction* is essentially an object that reserves a connection to the datastore. It lets you do all of the operations we’ve seen thus far, but guarantees that they’ll be **executed on the same connection**.

You begin a transaction with a call to `db.Begin()`, and close it with a `Commit()` or `Rollback()` method on the resulting `Tx` variable. Under the covers, the `Tx` gets a connection from the pool, and reserves it for use only with that transaction. The methods on the `Tx` map one-for-one to methods you can call on the database itself, such as `Query()` and so forth.  
Prepared statements that are created in a transaction are bound exclusively to that transaction.

Don't not mingle the use of transaction-related functions such as `Begin()` and `Commit()` with SQL statements such as `BEGIN` and `COMMIT` in your SQL code.
- The `Tx` objects could remain open, reserving a connection from the pool and not returning it.
- The state of the database could get **out of sync** with the state of the Go variables representing it.
- You could believe you’re executing queries on a single connection, inside of a transaction, when in reality Go has created several connections for you invisibly and some statements aren’t part of the transaction.

While you are working inside a transaction you should be careful not to make calls to the `db` variable. Make all of your calls to the `Tx` variable that you created with `db.Begin()`. `db` is not in a transaction, only the `Tx` object is. If you make further calls to `db.Exec()` or similar, those will happen outside the scope of your transaction, on other connections.

If you need to work with multiple statements that modify connection state, you need a `Tx` even if you don’t want a transaction per se. For example:
- Creating temporary tables, which are only visible to one connection.
- Setting variables, such as MySQL’s `SET @var := somevalue` syntax.
- Changing connection options, such as character sets or timeouts.

If you need to do any of these things, you need to bind your activity to a single connection, and the only way to do that in Go is to use a `Tx`.
# Using Prepared Statements
Prepared statements have all the usual benefits in Go: security, efficiency, convenience. 

At the **database level**, a prepared statement is bound to **a single database connection**.  
The **typical flow** is that the client sends a SQL statement with placeholders to the server for preparation, the server responds with a statement ID, and then the client executes the statement by sending its ID and parameters.

In Go, however, connections are not exposed directly to the user of the `database/sql` package. You don’t prepare a statement on a connection. You prepare it on a `DB` or a `Tx`. And `database/sql` has some convenience behaviors such as automatic retries. For these reasons, the underlying association between prepared statements and connections, which exists at the driver level, is hidden from your code. **Workflow**:
1. When you prepare a statement, it’s prepared on a connection in the pool.
2. The `Stmt` object remembers which connection was used.
3. When you execute the `Stmt`, it tries to use the connection. If it’s not available because it’s closed or busy doing something else, it gets another connection from the pool and re-prepares the statement with the database on another connection.

Because statements will be **re-prepared** as needed when their original connection is busy, it’s possible for high-concurrency usage of the database, which may keep a lot of connections busy, to create a large number of prepared statements. This can result in apparent leaks of statements, statements being prepared and re-prepared more often than you think, and even running into server-side limits on the number of statements.

Go creates prepared statements for you under the covers. A simple `db.Query(sql, param1, param2)`, for example, works by preparing the sql, then executing it with the parameters and finally closing the statement.

Sometimes a prepared statement is not what you want:
1. The database doesn’t support prepared statements. When using the MySQL driver, for example, you can connect to MemSQL and Sphinx, because they support the MySQL wire protocol. But they don’t support the “binary” protocol that includes prepared statements, so they can fail in confusing ways.
2. The statements aren’t reused enough to make them worthwhile, and security issues are handled in other ways, so performance overhead is undesired.  
- [x] [An example](https://www.vividcortex.com/blog/2014/11/19/analyzing-prepared-statement-performance-with-vividcortex/)

If you don’t want to use a prepared statement, you need to use `fmt.Sprint()` or similar to assemble the SQL, and pass this as the only argument to `db.Query()` or `db.QueryRow()`. And your driver needs to support plaintext query execution, which is added in Go 1.1 via the `Execer` and `Queryer` interfaces, [documented here](http://golang.org/pkg/database/sql/driver/#Execer).

Prepared statements that are created in a `Tx` are bound exclusively to it, so the earlier cautions about re-preparing do not apply. When you operate on a `Tx` object, your actions map directly to the one and only one connection underlying it.  
This also means that prepared statements created inside a `Tx` can’t be used separately from it. Likewise, prepared statements created on a `DB` can’t be used within a transaction, because they will be bound to a different connection.

To use a prepared statement prepared outside the transaction in a `Tx`, you can use `Tx.Stmt()`. It does this by taking an existing prepared statement, setting the connection to that of the transaction and re-preparing all statements every time they are executed. (This behavior and its implementation are undesirable and there’s even a TODO in the `database/sql` source code to improve it; we advise against using this.)

[database/sql: Strange Errors when Closing a Tx's Prepared Statement after Commit](https://github.com/golang/go/issues/4459)
# Handling Errors

```go
for rows.Next() {
	// ...
}
if err = rows.Err(); err != nil {
	// handle the error here
}
```

The error from `rows.Err()` could be the result of a variety of errors in the `rows.Next()` loop. The loop might exit for some reason other than finishing the loop normally, so you always need to check whether the loop terminated normally or not. An abnormal termination automatically calls `rows.Close()`, although it’s harmless to call it multiple times.

```go
for rows.Next() {
	// ...
	break // whoops, rows is not closed! memory leak...
}

// do the usual "if err = rows.Err()" [omitted here]

// it's always safe to [re?]close here
if err = rows.Close(); err != nil {
	log.Println(err)
}
```

You should always explicitly close a `sql.Rows` if you exit the loop prematurely, as previously mentioned. It’s auto-closed if the loop exits normally or through an error.  
The error returned by `rows.Close()` is the only exception to the general rule that it’s best to capture and check for errors in all database operations. If `rows.Close()` returns an error, it’s unclear what you should do. Logging the error message or panicing might be the only sensible thing, and if that’s not sensible, then perhaps you should just ignore the error.

```go
var name string
err = db.QueryRow("select name from users where id = $1", 1).Scan(&name)
if err != nil {
	log.Fatal(err)
}
fmt.Println(name)

var name string
err = db.QueryRow("select name from users where id = $!", 1).Scan(&name)
if err != nil {
	if err == sql.ErrNoRows {
		// there were no rows, but otherwise no error occurred
	} else {
		log.Fatal(err)
	}
}
fmt.Println(name)

```

If there was no user with `id = 1`, then there would be no row in the result, and `Scan()` would not scan a value into `name`. Go defines a special error constant, called `sql.ErrNoRows`, which is returned from `QueryRow()` when the result is empty. This needs to be handled as a special case in most circumstances.  
You should only run into this error when you’re using `QueryRow()`. If you encounter this error elsewhere, you’re doing something wrong.

If your connection to the database is dropped, killed, or has an error, you don’t need to implement any logic to retry failed statements. As part of the connection pooling in `database/sql`, handling failed connections is built-in. If you execute a query or other statement and the underlying connection has a failure, Go will reopen a new connection (or just get another from the connection pool) and retry, up to 10 times.
# Working with NULLs
Nullable columns are annoying and lead to a lot of ugly code. If you can, avoid them. If not, then you’ll need to use special types from the `database/sql` package to handle them, or define your own.

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

If you need to define your own types to handle NULLs, you can copy the design of `sql.NullString` to achieve that.

Nullability can be tricky, and not future-proof. If you think something won’t be null, but you’re wrong, your program will crash, perhaps rarely enough that you won’t catch errors before you ship them.  
One of the nice things about Go is having a useful default zero-value for every variable. This isn’t the way nullable things work.

```go
rows, err := db.Query(`
	SELECT
		name,
		COALESCE(other_field, '') as otherField
	WHERE id = $1
`, 42)

for rows.Next() {
	err := rows.Scan(&name, &otherField)
	// ..
	// If `other_field` was NULL, `otherField` is now an empty string. This works with other data types as well.
}
```

If you can’t avoid having `NULL` values in your database, there is another work around that most database systems support, namely `COALESCE()`（返回可变参数中左侧开始第 1 个不是 `NULL` 的值）.
# Working with Unknown Columns

```go
// cols: []string
cols, err := rows.Columns()
if err != nil {
	// handle the error
} else {
	dest := []interface{}{ // Standard MySQL columns (8)
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
	} else if len(cols) > 8 {
		// Handle this case
	}
	err = rows.Scan(dest...)
	// Work with the values in dest
}
```

If you don’t know how many columns the query will return, you can use `Columns()` to find a list of column names. You can examine the length of this list to see how many columns there are, and you can pass a slice into `Scan()` with the correct number of values.  
Some forks of MySQL return different columns for the `SHOW PROCESSLIST` command, so you have to be prepared for that or you’ll cause an error.

```go
cols, err := rows.Columns() // Remember to check err afterwards
vals := make([]interface{}, len(cols))
for i, _ := range cols {
	vals[i] = new(sql.RawBytes)
}
for rows.Next() {
	err = rows.Scan(vals...)
	// Now you can check each element of vals for nil-ness,
	// and you can use type introspection (reflection) and type assertions
	// to fetch the column into a typed variable.
}
```

If you don’t know the columns or their types, you should use `sql.RawBytes`.
# The Connection Pool
Connection pooling means that executing two consecutive statements on a single database might open two connections and execute them separately. For example, `LOCK TABLES` followed by an `INSERT` can block because the `INSERT` is on a connection that does not hold the table lock.

Connections are created when needed and there isn’t a free connection in the pool.

By default, there’s no limit on the number of connections. If you try to do a lot of things at once, you can create an arbitrary number of connections. This can cause the database to return an error such as “too many connections.”

In Go 1.1 or newer, you can use `db.SetMaxIdleConns(N)` to limit the number of idle connections in the pool. This doesn’t limit the pool size, though.  
Connections are recycled rather fast. Setting a high number of idle connections with `db.SetMaxIdleConns(N)` can reduce this churn, and help keep connections around for reuse.  
Keeping a connection idle for a long time can cause problems. Try db.`SetMaxIdleConns(0)` if you get connection timeouts because a connection is idle for too long.

In Go 1.2.1 or newer, you can use `db.SetMaxOpenConns(N)` to limit the number of total open connections to the database.

You can also specify the maximum amount of time a connection may be reused by setting `db.SetConnMaxLifetime(duration)` since reusing long lived connections may cause network issues. This closes the unused connections lazily i.e. closing expired connection may be deferred.
# Surprises, Anti-patterns and Limitations
## Resource Exhaustion
Opening and closing databases can cause exhaustion of resources.

Failing to (read all rows or use `rows.Close()`) reserves connections from the pool.

Using `Query()` for a statement that doesn’t return rows will reserve a connection from the pool.

Failing to be aware of how prepared statements work can lead to a lot of extra database activity.

## Large `uint64` Values

```go
_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64) // Error
```

You can’t pass big unsigned integers as parameters to statements if their high bit is set.  
Be careful if you use `uint64` values, as they may start out small and work without error, but increment over time and start throwing errors.
## Connection State Mismatch
Some connection state, such as whether you’re in a transaction, should be handled through the Go types instead.  
For example, setting the current database with a `USE` statement is a typical thing for many people to do. But in Go, it will affect only the connection that you run it in. Unless you are in a transaction, other statements that you think are executed on that connection may actually run on different connections gotten from the pool, so they won’t see the effects of such changes. Additionally, after you’ve changed the connection, it’ll return to the pool and potentially **pollute** the state for some other code. This is one of the reasons why you should never issue `BEGIN` or `COMMIT` statements as SQL commands directly, too.
## Multiple Result Sets
The Go driver doesn’t support multiple result sets from a single query in any way, and there doesn’t seem to be any plan to do that, although there is a [feature request](https://github.com/golang/go/issues/5171) for supporting bulk operations such as bulk copy. This means, among other things, that a stored procedure that returns multiple result sets will not work correctly.
## Invoking Stored Procedures
## Multiple Statement Support

```go
_, err := db.Exec("DELETE FROM tbl1; DELETE FROM tbl2") // Error/unpredictable result
```

The `database/sql` doesn’t explicitly have multiple statement support, which means that the behavior of this is backend dependent. The server is allowed to interpret this however it wants, which can include returning an error, executing only the first statement, or executing both.

```go
rows, err := db.Query("select * from tbl1") // Uses connection 1
for rows.Next() {
	err = rows.Scan(&myvariable)
	// The following line will NOT use connection 1, which is already in-use
	db.Query("select * from tbl2 where id = ?", myvariable)
}
```

When you’re not working with a transaction, it is perfectly possible to execute a query, loop over the rows, and within the loop make a query to the database (which will happen on a new connection).

```go
tx, err := db.Begin()
rows, err := tx.Query("select * from tbl1") // Uses tx's connection
for rows.Next() {
	err = rows.Scan(&myvariable)
	// ERROR! tx's connection is already busy!
	tx.Query("select * from tbl2 where id = ?", myvariable)
}
```

There is no way to batch statements in a transaction. Each statement in a transaction must be executed serially, and the resources in the results, such as a Row or Rows, must be scanned or closed so the underlying connection is free for the next statement to use.  
Go doesn’t stop you from trying, though. For that reason, you may wind up with a corrupted connection if you attempt to perform another statement before the first has released its resources and cleaned up after itself. This also means that each statement in a transaction results in a separate set of network round-trips to the database.
