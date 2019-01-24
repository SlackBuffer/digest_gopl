- Conventions
    1. Describe each package in a comment immediately preceding its package declaration
        - For a `main` package, this comment is one or more complete sentences that describe the program as a while
    2. Write a comment before the declaration of each function to specify its behavior
    - > These conventions are used by tools like `go doc` and `godoc` to locate and display documentation
    - Short names are preferred, especially for local variables
        - The larger the scope of a name, the longer and more meaningful it should be
    - Use "camel case"
    - The letters of acronyms and initialisms like ASCII and HTML are always rendered in the **same case**
- Go does not permit unused local variables. Use the **blank identifier `_`** whenever syntax requires a variable name but program logic does not
# Good practices
- `if`, `switch`, `for` statements can include an optional simple statement-a short variable declaration, an increment or assignment statement, or a function call-that can be used to set a value before it is tested

    ```go
    // good
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    // verbose.
    err := r.ParseForm()
    if err != nil {
        log.Print(err)
    }
    ```