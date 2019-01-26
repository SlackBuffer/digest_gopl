- Function and other package-level entities may be declared in any order
- An entity declared within a function is local to that function
- An entity declared outside of a function is visible in all files of the package to which it belongs (package-level entity), as if the source code were all in a single file
- A entity whose name begins with an upper-case letter (**exported**)is visible and accessible outside its own package and may be referred to by other parts of the program
- The lifetime of a package-level variable is the entire execution of the program
- A new instance of a local variable is created each time the declaration is executed, and lives on util it becomes unreachable, at which point its storage may be recycled
    - Every package-level variable, and every local variable of each currently active function, can potentially be the start or root of a path to the variable in question, (variable) following pointers and other kinds of reference that ultimately lead to the variable
    - If no such path exists, the variable has become unreachable, so it can no longer affect the rest of the computation
- A compiler may choose to allocate local variables on the heap or on the stack. The choice is not determined by whether `var` or `new` was used to declare the variable

    ```go
    var global *int
    func f() {
        var x int
        x = 1
        global = &x
    }
    func g() {
        y := new(int)
        *y = 1
    }
    ```

    - `x` here must be heap-allocated because it's still reachable from `global` after `f` (function context exists on the stack) has returned, despite being declared as a local variable; we say `x` escapes from `f`
    - When `g` returns, the variable `*y` becomes unreachable and can be recycled. Since `*y` does not escape from `g`, it's safe for the compiler to allocate `*y` on the stack, even though it was allocated with `new`