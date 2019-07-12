- Tuple assignment
    - **All of the right-hand side expressions are evaluated before any of the variables are updated**, making this form most useful when some of the variables appear on both sides of the assignment

    ```go
    i, j = j, i // swap values
    // greatest common divisor
    func gcd(x, y int) int {
        for y != 0 {
            x, y = y, x%y
        }
        return x
    }
    // n-th Fibonacci number
    func fib(n int) int {
        x, y := 0, 1
        for i := 0; i < n; i++ {
            x, y = y, x+y
        }
        return x
    }

    f, err = os.Open("foo.txt")
    v, ok = m[key]  // map lookup
    v, ok = x.(T)   // type assertion
    v, ok = <-ch    // channel receive
    ```

- Assignment statements are an explicit form of assignment
- Places where an assignment occurs implicitly
    1. A function call implicitly assigns the argument values to the corresponding parameter variables
    2. A `return` statement implicitly assigns the `return` operands to the corresponding result variables
    3. A literal expression for a composite type (slices, maps, channels)
    
        ```go
        medals := []string{"gold", "silver", "bronze"}
        // implicitly assigns each element
        medals[0] = "gold" // ...
        ```
