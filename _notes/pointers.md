- A pointer value is the address of a variable
    - Not every value has an address, but every variable does
- A pointer is the location at which a value is stored
- Pointer can be used to read or update the value of a variable indirectly, without using or even knowing the name of the variable, if indeed it has a name (`new(T)` creates an unnamed variable)
- `var x int`
    - The expression `&x` yields a pointer to an integer variable, that is, **a type of `*int`** (pronounced "pointer to int")

    ```go
    x := 1
    p := &x // p, of type *int, points to x
    ```

- Each component of a variable of aggregate type-a field of a struct or an element of an array-is also a variable and thus has an address too
- Variables are also described as **addressable values**. Expressions that denote variables are the only expressions to which the address-of operator `&` may be applied
- Pointers in Go are explicitly visible, but there's no pointer arithmetic
- `&` yields the address of a variable; `*` retrieves the variable that the pointer refers to
- The zero value for a pointer of any type is `nil`
- `p != nil` is true if `p` points to a variable
- Pointers are comparable; 2 pointer are equal if and only if they point to the same variable or both are `nil`
- Perfectly safe for a function to return the address of a local variable

    ```go
    var p = f()
    func f() *int {
        v := 1
        return &v
    }
    fmt.Println(f() == f()) // "false"
    ```

    - Local variable `v` remains in existence even after the call has returned. The pointer `p` will still refer to it
    - Each call of `f` returns a distinct value
- Each time we take the address of a variable or copy a pointer, we create new aliases to identify the same variable
- Aliasing also occurs when we copy values of other reference types like slices, maps, and channels, and even structs, arrays, and interfaces that contain these types
- Pointer aliasing is useful but is a double-edged sword: to find all the statements that access a variable, we have to know all its aliases