- A channel is a **communication mechanism** that allows one goroutine to pass values of a specified type to another goroutine
- `main` runs in a goroutine and the `go` statement creates additional goroutines
- When one goroutine attempts a send or receive **on a channel**, it blocks until another goroutine attempts the corresponding receive or send operation, at which point the value is transferred and both goroutines proceed
- Communicating sequential process (CSP) - a model of concurrency in which values are passed between independent activities (goroutines) but variables are for most part confined to a single activity
# Goroutines
- In Go, each concurrently executing activity is called a goroutine
- When a program starts, its only goroutine is the one that calls the `main` function (called **main goroutine**)
- New goroutine are created by the `go` statement (new goroutine calls other functions)
    - Syntactically, a `go` statement is an ordinary function or method call prefixed by the keyword `go`
    - A `go` statement causes the function to be called in a newly created goroutine
- The `go` keyword itself completes immediately

    ```go
    f() // wait for f() to return
    go f() // create a new goroutine that calls f(); won't wait
    ```

- When the `main` function returns, **all** goroutines are abruptly terminated and the **program exists**
    - Other than by returning from `main` or existing the program, there is no programmatic way for one goroutine to stop another
    - There are ways to communicate with a goroutine to request that it stop itself
## Example: concurrent clock server
- `time.Time.Format` provides a way to format date and time information by example
    - Its argument is a template indicating how to format a reference time, specially `Mon Jan 2 03:04:05PM 2006 UTC-0700`
    - The reference time has 8 components. Any collection of them can appear in the `Format` string in **any order** and in **a number of formats**
    - The selected components of will be displayed in the selected formats
    - The `time` package defines templates for many standard time formats， such as `time.RFC1123`
    - The same mechanism is used in reverse when passing a time using `time.Parse`