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