- Conventions
    - Comments
        1. Describe each package in a comment immediately preceding its package declaration
           - The *doc comment* immediately preceding the `package` declaration documents the package as a whole
            - For a `main` package, this comment is one or more complete sentences that describe the program as a while
           - Only **one file** in each package should have a package doc comment
           - Extensive doc comments are often placed in a file of their own, conventionally called `doc.go`
        2. Write a comment before the declaration of each function to specify its behavior
        - These conventions are used by tools like `go doc` and `godoc` to locate and display documentation
    - Short names are preferred, especially for local variables. The larger the scope of a name, the longer and more meaningful it should be
    - Use "camel case"
    - The letters of acronyms and initialisms like ASCII and HTML are always rendered in the **same case**
- Go does not permit unused local variables. Use the **blank identifier `_`** whenever syntax requires a variable name but program logic does not
- Simplify redundant boolean expressions like `x==true` to `x`
- Functions that merely access or modify internal values of a type, such as the methods of the `Logger` type from `log` package, are called *getter* and *setter*
    - When naming a getter method, we usually **omit** the `Get` prefix. This preference for brevity extends to all methods, not just field accessors, and to other redundant prefixes as well, such as `Fetch`, `Find`, and `Lookup`
- By convention, the variables guarded by a mutex are declared immediately after the declaration of the mutex itself
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

- Deal with the error in the `if` block and then return, so that the successful execution path is not indented

    ```go
    // normal practice
    f, err := os.Open(fname)
    if err != nil {
        return err
    }
    f.Stat()
    f.Close()
    // don't
    if f, err := os.Open(name); err != nil {
        return err
    } else {
        f.Stat()
        f.Close()
    }
    ```