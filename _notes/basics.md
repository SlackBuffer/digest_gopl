- Doesn't require **semicolons** at the ends of statements or declarations
    - Except where 2 or more appear on the same line
    - Newlines following certain tokens are converted into semicolons
        1. The opening brace `{` of the function must be on the same line as the end of the `func` declaration, not on a line by itself
        2. In the expression `x + y`, a newline is permitted after but not before the `+` operator
- All **indexing** in Go uses **half-open** intervals that include the first index but exclude the last
- `++`, `--`
    - **Prefix** only
    - **Statements**, not expressions (so `j = i++` is illegal)
- Declaring a string variable

    ```go
    /* generally, use the first 2 forms */
    s := ""             // used only within function, not for package-level variables
    var s string        // relies on default initialization
    var s = ""          // rarely used except when declaring multiple variables
    var s string = ""   // explicit about the variable's type, necessary when it is not the same as that of the initial value
    ```