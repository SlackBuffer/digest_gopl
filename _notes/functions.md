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