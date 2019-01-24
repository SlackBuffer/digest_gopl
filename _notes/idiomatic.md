- The `import` declarations must follow the `package` declaration
- Conventions
    - Describe each package in a comment immediately preceding its package declaration
        - For a `main` package, this comment is one or more complete sentences that describe the program as a while
    - Write a comment before the declaration of each function to specify its behavior
    - > These conventions are used by tools like `go doc` and `godoc` to locate and display documentation
- Go does not permit unused local variables. Use the **blank identifier `_`** whenever syntax requires a variable name but program logic does not
# Good practices
- `if`, `switch`, `for` statements can include an optional simple statement-a short variable declaration, an increment or assignment statement, or a function call-that can be used to set a value before it is tested

    ```go
    // good
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    // longer...
    err := r.ParseForm()
    if err != nil {
        log.Print(err)
    }
    ```