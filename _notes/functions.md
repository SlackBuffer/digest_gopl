# Declarations
- A function declaration has a name, an optional list of parameters, an optional list of results, and a body

    ```go
    func name(parameter-list) (result-list) {
        body
    }
    ```

    - Parameters are local variables whose value or **arguments** are supplied by the caller
    - If the function returns one unnamed result or no result at all, parentheses are optional and usually omitted
    - Results may be named. Each name declares a local variable initialized to the zero value for its type
- A function that has a result list must end with a `return` statement unless execution clearly cannot reach the end of the function, perhaps because the function ends with a call to `panic` or an infinite `for` loop with no `break`
- A sequence of parameters or results of the same type can be factored so that the type is written only once
- 4 ways to declare a function with 2 parameters and 1 result, all of type `int`

    ```go
    func add(x int, y int) int      { return x + y }
    func sub(x, y int) (z int)      { z = x - y; return }
    func first(x int, _ int) int    { return x }
    func zero(int, int) int         { return 0 }
    ```

    - The blank identifier can be used to emphasize that a parameter is unused
- The type of a function is called its **signature**
- 2 functions have the same type or signature if they have the same sequence of parameter types and the same sequence of result types
    - The names of parameters and results don't affect the type, nor does whether or not they were declared using the factored form
- Every function call must provide an argument for each parameter, in the order in which the parameters were declared
- Go has no concept of default parameter values, nor any way to specify arguments by name, so the names of parameters and results don't matter to the caller except as documentation
- Parameters are local variables within the body of the function, with their values set to the arguments supplied by the caller
- Function parameters and named results are variables in the same lexical block as the function's **outermost local variables**
- Function parameters and results are created each time their enclosing function is called
- Execution of the function begins with the first statement and continues until it encounters a `return` statement or reaches the end of a function that has no results. Control and any results are then returned to the caller
- A function declaration without a body indicates that the function is implemented in a language other than Go

    ```go
    package math
    func Sin(x float64) float64 // implemented in assembly language
    ```

# Recursion
- Recursive functions may call themselves directly or indirectly
- Many programming language implementations use a fixed-size function call stack; sizes from 64KB to 2MB are typical
- Fixed-size stacks impose a limit on the depth of recursion, so one must be careful to avoid a stack overflow when traversing large data structures recursively; fixed-size stacks may even pose a security risk
- Typical Go implementations use variable-size stacks that starts small and grow as needed up to a limit on the order of a gigabyte
    - This lets us use recursion safely without worrying about overflow
# Multiple return values
- The result of a multi-valued call may itself be returned from a (multi-valued) calling function
- A multi-valued call may appear as the sole argument when calling a function of multiple parameters
    - Rarely used in production code
    - Convenient during debugging since it lets us print all the result of a call using a single statement

    ```go
    log.Println(findLinks(url))

    links, err := findLinks(url)
    log.Println(links, err)
    ```

- Well-chosen names can document the significance of a function's results
- It's not always necessary to name multiple results solely for documentation
    - **Convention** dictates that a final `bool` result indicates success; an `error` result often needs no explanation
- In a function with named results, the **operands of a return statement may omitted**, called a **bare return**
- A bare return is shorthand way to return each of the named result variables in order
- In functions with many return statements and several results, bare returns can reduce code duplication, but they rarely make code  easier to understand
    - Bare returns are best used sparingly
# Errors
- Some functions always succeed at their task
    - `strings.Contains` and `strconv.FormatBool` have well-defined results for all possible argument values and cannot fail
- Other functions always succeed as long as their preconditions are met
- For many other functions, even in well-written program, success is not assured because it depends on factors beyond the programmer's control
    - Any functions that does I/O, for example, must confront the possibility of error
    - Indeed, it's when the most reliable operations fail unexpectedly that we most need know why
- Errors are thus an important part of a package's API or an application's user interface, and failure is just one of several expected behaviors
- A function for which failure is an expected behavior returns an additional result, conventionally the last one
    - If the failure has only one possible cause, the result is a boolean, usually `ok`
    - More often, and especially for I/O, the failure may have a variety of causes for which the caller need an explanation. In such cases, the type of the additional result is `error`
- The built-in type `error` is an interface type
    - An `error` may be nil or non-nil, that `nil` implies success and non-nil implies failure
    - A non-nil `error` has an error message string that can be obtained by calling its `Error` method or print by calling `fmt.Println(err)` or `fmt.Printf("%v", err)`
- Usually when a function returns a non-nil error, its other results are undefined and should be **ignored**
- However, a few functions may return partial results in error case
    - If an error occurs while reading from a file, a call to `Read` returns the number of bytes it was able to read and an `error` value describing the problem
    - For correct behavior, some callers may need to process the incomplete data before handling the error, so it's important that such functions clearly document their results
- Although Go does have an exception mechanism, it's used only for reporting truly **unexpected errors** that indicate a bug, not the **routine errors** that a robust program should be built to expect
    - The reason for this design is that exceptions tend to entangle the **description of an error** with the **control flow** required to handle it, often leading to an undesirable outcome: routine errors are reported to the end user in the form of an incomprehensible stack trace, full of information about the structure of the program but lacking intelligible context about what went wrong
- Go programs use ordinary control-flow mechanisms like `if` and `return` to respond to errors
    - This style demands that more attention be paid to error-handling logic, but this is precisely the point
## Error-handling strategies
- When a function call returns an error, it's the caller's responsibility to check it and take appropriate action
- After checking an error, failure is usually dealt with **before success**
- If failure causes the function to return, the logic for success is not indented with `else` block but follows at the outer level
- Functions tend to exhibit a common structure, with a series of initial checks to reject errors, followed by the substance of the function at the end, **minimally indented**
### 1. Propagate the error
- The most common is to propagate the error, so that a failure in a subroutine becomes a failure of the calling routine
- `fmt.Errorf` formats an error message using `fmt.Sprintf` and returns a new `error` value
    - It's used to build descriptive errors by **successively prefixing** additional context information to the original error message

    ```go
    doc, err := html.Parse(resp.Body)
    resp.Body.Close()
    if err != nil {
        return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
    }

    resp, err := http.Get(url) 
    if err != nil {
        return nil, err
    }
    ```

- When the error is handled by the program's `main` function, it should provide a clear causal chain from the root problem to the overall failure
    - `genesis: crashed: no parachute: G-switch failed: bad relay orientation`
- Because error messages are frequently chained together, message strings should not be capitalized and **newlines should be avoided**
    - The resulting errors may be long, but they'll be self-contained when found by tools like `grep`
- When designing error messages, be deliberate, so that each one is a meaningful description of the problem with sufficient and relevant detail, and be consistent, so that errors returned by the same function or by a group of functions in the same packages are similar in form and can be dealt with in the same way
    - The `os` package guarantees that every error returned by a file operation describes not just the nature of the failure (permission denied, no such directories, and so on) but also the name of the file, so the caller needn't include this information in the error message it constructs
- In general, the call `f(x)` is responsible for reporting the attempted operation `f` and the argument value `x` as they relate to the context of the error. The caller is responsible for adding further information that it has but the call `f(x)` does not
### 2. Retry
- For errors that represent transient or unpredictable problems, it may make sense to retry the failed operation, possibly with a delay between tries, and perhaps with a limit on the number of attempts or the time spent trying before giving up entirely
### 3. Print the error and stop the program
- If progress is impossible, the caller can print the error and stop the program gracefully
- This course of action should generally be reserved for the **main package** of a program
    - Library functions should usually propagate errors to the caller, unless the error is a sign of an internal inconsistency - that is, a bug

    ```go
    // in function main
    if err := WaitForSever(url); err != nil {
        fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
        os.Exit(1)
    }
    ```

- A more convenient way to achieve the same effect is to call `log.Fatalf`. As with all the `log` functions, by default it prefixes the time and date to the error message

    ```go
    if err := WaitForServer(url); err != nil {
        log.Fatalf("Site is down: %v\n", err)
    }
    ```

    - For a more attractive output, we can set the prefix used by the `log` package to the name of the command, and suppress the display of display of the date and time

        ```go
        log.SetPrefix("wait: ")
        log.SetFlags(0)
        ```

### 4. Log and continue
- In some cases, it's sufficient just to log the error and then continue

    ```go
    if err := Ping(); err != nil {
        log.Printf("ping failed: %v; networking disabled", err)    
    }

    if err := Ping(); err != nil {
        fmt.Fprintf(os.Stderr, "ping failed: %v; networking disabled\n", err)
    }
    ```

    - All `log` functions append a newline if one is not already present
### 5.Ignore
- In rare cases we can safely ignore an error entirely

    ```go
    dir, err := ioutil.TempDir("", "scratch")
    if err != nil {
        return fmt.Errorf("failed to create temp dir: %v", err)
    }
    // ...use temp dir...
    os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodically
    ```

    - The call to `os.RemoveAll` may fail, but the program ignores it because the operating system periodically cleans out the temporary directory
    - In this case, discarding the error was intentional, but the program logic would be the same had we forget to deal with it
- Get into the habit of considering errors after every call, and when you deliberately ignore one, document your intent clearly
## End Of File (EOF)
- On occasion, a program must take different actions depending on the kind of error that has occurred
- Consider an attempt to read `n` bytes of data from a file
    1. If `n` is chosen to be the length of the file, any error represents a failure
    2. If the caller repeatedly tries to read fixed-size chunks until the file is exhausted, the caller must respond differently to an end-of-file condition that it does to all other errors
- The `io` package guarantees that any read failure caused by an end-of-file condition is always reported by a distinguished error, `io.EOF`

    ```go
    package io
    import "errors"
    // EOF is the error returned by Read when no more input is available
    var EOF = errors.New("EOF")
    ```

- The caller can detect this condition using a simple comparison

    ```go
    in := bufio.NewReader(os.Stdin)
    for {
        r, _, err := in.ReadRune()
        if err == io.EOF {
            break // finished reading
        }
        if err != nil {
            return fmt.Errorf("read failed: %v", err)
        }
        // ...use r...
    }
    ```

- Since in an end-of-file condition there's no information to report besides the fact of it, `io.EOF` has a fixed error message, "EOF"
    - For other errors, we may need to report both the quantity and quantity of the error, so to speak, so a fixed error value will not do
# Function values
- Functions are first-class values in Go: like other values, function values have **types**, and they may be assigned to variables or passed to or returned from functions
- A function value may be called like any other function

    ```go
    func square(n int) int { return n * n }
    func negative(n int) int { return -n }
    func product(m, n int) int { return m * n }

    f := square
    fmt.Println(f(3)) // "9"
    f = negative
    fmt.Printf("%T\n", f) // "func(int) int"
    f = product // compile error: can't assign func(int, int) int to func(int) int
    ```

- The zero value of a function type is `nil`
- Calling a nil function value causes a panic

    ```go
    var f func(int) int
    f(3) // panic: call of nil function
    ```

- Function values may be compared with `nil`

    ```go
    var f func(int) int
    if f != nil {
        f(3)
    }
    ```

- Function values are not comparable, so they may not be compared against each other or used as keys in a map
- Function values let us parameterize our functions over not just data, but **behavior** too
    - `strings.Map` applies a function to each character of a string, joining the results to make another string

        ```go
        func add1(r rune) rune { return r + 1 }
        fmt.Println(strings.Map(add1, "HAL-9000")) // IBM.:111
        fmt.Println(strings.Map(add1, "VMS")) // "WNT"
        fmt.Println(strings.Map(add1, "Admin")) // "Benjy"
        ```
    
# Anonymous functions
- Named functions can be declared only at package level
- A function literal can be used to denote a function value within any expression
- A function literal is written like a function declaration, but without a name following the `func` keyword
- A function literal is an expression, and its **value** is called an anonymous function
- Function literals let us define a function at its point of use

    ```go
    strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")
    ```

- Functions defined in this way has access to **the entire lexical environment**, so the inner function can refer to variables from the enclosing function
- Function values can have state (closures)
- The anonymous inner function can access and update the local variables the local variables of the enclosing function
- These **hidden variable references** are why we classify functions as **reference types** and why function values are not **comparable**
- When an anonymous function requires recursion, we must first declare a variable, and then assign the anonymous function to that variable
    - Had these 2 steps been combined in the declaration, the function literal would not be within the scope of the variable (`visitAll`) (`:=` 类型推断) so it would have no way to call itself recursively

    ```go
    var visitAll func(items []string)
	visitAll = func(items []string) {
        // ...
        visitAll(m[item])
        // ...
    }
    
    visitAll := func(item []string) {
        // ...
        visitAll(m[item]) // compile error: undefined: visitAll
        // ...
    }
    ```

- Crawling the web is a problem of graph traversal
- DFS, BFS
    - https://www.youtube.com/watch?v=bIA8HEEUxZI
## Caveat: capturing iteration variables
- Consequence of the **scope rules for loop variables**

    ```go
    var rmdirs []func()
    for _, d := range tempDirs() {
        dir := d // necessary!

        os.MkdirAll(dir, 0755) // this creates necessary parent directories too
        rmdirs = append(rmdirs, func() {
            os.RemoveAll(dir)
        })
    }
    for _, rmdirs := range rmdirs {
        rmdir() //clean up
    }
    ```

    - The `for` loop introduces a new lexical block in which the variable `d` is declared
    - Without `dir := d`, all function values created by this loop "capture"" and share the same variable `d` - an **addressable storage location**, not its value at that particular moment. The value of `d` is updated in successive iterations
- Frequently, the inner variable introduced to work around this problem is given the exact name as outer variable of which it is a copy, leading to odd-looking but crucial variable declaration

    ```go
    for _, dir := range tempDirs() {
        dir := dir
    }
    ```

- The problem of iteration variable capture is most often encountered when using the `go` statement or with `defer` since both may delay the execution of a function value until after the hoop has finished
# Variadic functions
- A variadic function is one that can be called with varying numbers of arguments
    - The most familiar examples are `fmt.Printf` and its variants
        - `Printf` requires one fixed argument at the beginning, then accepts any number of subsequent arguments
- To declare a variadic function, the type of the final parameter is preceded by an ellipsis `...`, which indicates that the function may be called with any number of arguments of this type

    ```go
    func sum(vals ...int) int {
        total := 0
        for _, val := range vals { // vals is an []int
            total += val
        }
        return total
    }

    values := []int{1, 2, 3, 4}
    fmt.Println(sum(values...)) // 10
    ```

- Implicitly, the caller allocates an array, copies the arguments into it, and passes a slice of the entire array to the function
- Variadic functions are often used for string formatting
    - The suffix `f` is a widely followed naming convention for variadic functions that accept a `Printf`-style format string

    ```go
    // the interface{} means it can accept any values at all for its final arguments
    func errorf(linenum int, format string, args ...interface{}){
        fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
        fmt.Fprintf(os.Stderr, format, args...)
        fmt.Fprintln(os.Stderr) // for newline
    }
    linenum, name := 12, "count"
    errorf(linenum, "undefined: %s", name) // "Line 12 : undefined: count"
    ```

# Deferred function calls
- Using the output of `http.Get` as the input to `html.Parse` works well if the content of the required URL is HTML
- Many pages contains images, plain text, and other file formats. Feeding such files into an HTML parser could have undesirable effects
- The duplicated `resp.Body.Close()` call ensures that the network connection on all execution paths, including failures is closed. As functions grow more complex and have to handle more errors, such duplication of clean-up logic may become a maintenance problem
- A `defer` statement is an ordinary function or method call prefixed by the keyword `defer`
- The function and argument expressions are evaluated when the statement is executed, but the actual call is **deferred** util **the function that contains the `defer` statement has finished**, whether normally, by executing a return statement or falling off the end, or abnormally, by panicking
- Any number of calls may be deferred; they are executed in the **reverse of the order** in which they were deferred
- A `defer` statement is often used with paired operations like open and close, connect and disconnect, or lock and unlock to ensure that resources are released in all cases, no matter how complex the control flow
- The right place for a `defer` statement that releases a resource is **immediately after the resource has been successfully acquired**

    ```go
    func ReadFile(filename string) ([]byte, error) {
        f, err := os.Open(filename)
        if err != nil {
            return nil, err
        }
        defer f.Close()
        return ReadAll(f)
    }

    var mu sync.Mutex
    var m = make(map[string]int)
    func lookup(key string) int {
        mu.Lock()
        defer mu.Unlock()
        return m[key]
    }
    ```

- The `defer` function can also be used to **pair "on entry" and "on exit" actions** when debugging a complex function
    - By deferring a call to the returned function, we can instrument the entry point and all exit points of a function in a single statement and even pass value between the two actions (Trace)
- Deferred functions run **after** return statements have updated the function's result variables
- Because an anonymous function can access its enclosing function's variables, including named results, a deferred anonymous function can observe the function's results

    ```go
    func double(x, int) (result int) {
        defer func() { fmt.Printf("double(%d) = %d\n", x, result )}()
        return x + x
    }
    _ = double(4)
    ```

    - This trick may be useful in functions with many return statements
- A deferred anonymous function can change the values that the enclosing function returns to its caller

    ```go
    func triple(x int) (result int) {
        defer func() { result += x }()
        return double(x)
    }
    triple(4)
    ```

- Because deferred functions aren't executed until the very end of a function's execution, a `defer` statement in a loop deserves extra scrutiny

    ```go
    for _, filename := range filenames {
        f, err := os.Open(filename)
        if err != nil {
            return err
        }
        defer f.Close() // risky: could run out of file descriptors
        // ...process f...
    }
    ```

    - The code below could run out of file descriptors since no file will be closed until all files have been processed
    - One solution is to move the loop body, including the `defer` statement, into another function that's called on each iteration

        ```go
        for _, filename := range filenames {
            if err := doFile(filename); err != nil {
                return err
            }
        }
        func doFile(filename string) error {
            f, err := os.Open(filename)
            if err != nil {
                return err
            }
            defer f.Close()
            // ...process f...
        }
        ```

- Improved fetch

    ```go
    func fetch(url string) (filename string, n int64, err error) {
        resp, err := http.Get(url)
        if err != nil {
            return "", 0, err
        }

        defer resp.Body.Close()

        local := path.Base(resp.Request.URL.Path)
        if local == "/" {
            local = "index.html"
        }

        f, err := os.Create(local)
        if err != nil {
            return "", 0, err
        }

        n, err = io.Copy(f, resp.Body)
        if closeErr := f.Close(); err == nil {
            err = closeErr
        }
        return local, n, err
    }
    ```

    - It's tempting to use a second deferred call, to `f.Close()`, to close the local file, but this would be subtly wrong
    - `os.Create` opens a file for writing, creating it as needed. On many file systems, notably NFS, write errors are not reported immediately but may be **postponed** until the file is closed. Failure to check the result of the close operation could cause serious data loss to go unnoticed
    - If both `io.Copy` and `f.Close` fail, we should prefer to report the error from `io.Copy` since it occurred first and is more likely to tell us the root cause
# Panic
- Go's type system catches many mistakes at compile time, but others, like out-of-bounds array access or nil pointer dereference, requires checks at **run time**
- When the Go runtime detects these mistakes, it panics
- During a typical panic, **normal execution stops**, all deferred function calls **in that goroutine** are executed, and the program crashes with a log message
    - This log message includes the panic value, which is usually an error message of some sort, and, for each goroutine, a stack trace showing the stack of function calls that were active at the time of the panic
    - This log message often has enough information to diagnose the root cause of the problem without running the program again, so it should always be included in a bug report about a panicking program
- Not all panics come from the runtime
- The built-in `panic` function may be called directly; it accepts any value as an argument
- A panic is often the best thing to do when some impossible situation happens, for instance, execution reaches a case that logically can't happen
- Unless you can provide a more informative error message or detect an error sooner, there's no point asserting a condition that the runtime will check for you

    ```go
    func Reset(x *Buffer) {
        if x == nil {
            panic("x is nil") // unnecessary
        }
        x.elements = nil // this will panic automatically
    }
    ```

- Since a panic cause the program to crash, it's generally used for grave errors, such as a logical inconsistency in the program
- In a robust program, "expected" errors, the kind that arises from incorrect input, misconfiguration, or failing I/O, should be handled gracefully; they're best dealt with using `error` values
- `regexp.Compile` compiles a regular expression into an efficient form for matching
    - It returns an `error` if called with an ill-formed pattern, but checking this error is unnecessary and burdensome if the caller knows that a particular call cannot fail. In such cases, it's reasonable for the caller to handle an error by panicking, since it's believed to be impossible
    - Since most regular expressions are literals in the program source code, the `regexp` package provides a wrapper function `regexp.MustCompile` that does this check

    ```go
    package regexp
    func Compile(expre string) (*Regexp, error) { /* ... */ }
    func MustCompile(expr string) *Regexp {
        re, err := Compile(expr)
        if err != nil {
            panic(err)
        }
        return re
    }
    ```

    - The wrapper function makes it convenient for clients to initialize a package-level variable with a compiled regular expression

        ```go
        var httpSchemeRE = regexp.MustCompile(`^https?:`)
        ```
        
- Of course, `MustCompile` should not be called with untrusted input values
- The `Must` prefix is a common naming convention functions of this kind
- When a panic occurs, all deferred functions are run in **reverse order**, starting with those of the topmost function on the stack and proceeding up to the `main`
- For diagnostic purposes, the `runtime` package lets the programmer dump the stack using the same machinery (Defer2)
    - `runtime.Stack` can print information about functions that seems to have already been "unwound"
- Go's panic mechanism runs the deferred functions before it **unwind** (释放) the stack
# Recover
- Giving up is usually the right response to a panic, but not always
    - It might be possible to recover in some way, or at least clean up the mess before quitting
    - > A web server that encounters an unexpected problem could close the connection rather leave the client hanging, and during development, it might report the error to the client too
- If the built-in `recover` function is called within a deferred function and the function containing the `defer` statement is panicking, **`recover`** ends the current state of panic and **returns the panic value**
- The function that was panicking does not continue where it left off but **straightly returns normally** (with or without explicit `return` statement)
- If `recover` is called at any other time, it has no effect and returns `nil`
- The process of developing a parser for a language

    ```go
    func Parse(input string) (s *Syntax, err error) {
        defer func() {
            if p := recover(); p != nil {
                err = fmt.Errorf("internal error: %v", p)
            }
        }()
        // ...parser...
    }
    ```

    - It's preferred for the parser to turn panics into ordinary parser errors, perhaps with an extra message exhorting the user to file a bug report, instead of crashing
    - A fancier version might include the entire call stack using `runtime.Stack`
    - The deferred function then assigns to the `err` result, which is returned to the caller
- Recovering indiscriminately from panics is a dubious practice because the state of a package's variables after panic is rarely well defined or documented
- Further more, by replacing a crash with, say, a line in a log file, indiscriminate recovery may cause bugs to go unnoticed
- Recovering from a panic within the same package can help simplify the handling of complex or unexpected errors, but as a general rule, you **should not attempt to recover from another package's panic**
- Public APIs should report failures as `error`s
- Similarly, you should not recover from a panic that may pass through a function you do not maintain, such as a caller-provided callback, since you cannot reason about its safety
    - The `net/http` package provides a web server that dispatches incoming requests to user-provided handler functions
    - Rather than let a panic in one of these handlers kill the process, the server calls `recover`, prints a stack trace, and continue serving
    - This is convenient in practice, but it does risk leaking resources or leaving the failed handler in an unspecified state that could lead to other problems
- For all the reasons, it's safest to **recover selectively if at all**
    - Recover only from panic that were intended to be recovered from, which should be rare
        - This intention can be encoded by using a distinct, unexported type for the panic value and testing whether the value returned by `recover` has that type. If so, we report the panic as an ordinary `error`; if not, we call `panic` with the same value to resume the state of panic (Title3)
- For some conditions there is no recovery
    - Running out of memory, for example, causes the Go runtime to terminate the program with a fatal error