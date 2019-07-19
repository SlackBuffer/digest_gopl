<!-- - `if`
    - p8 -->
- `for`

    ```go
    for initialization; condition; post { /* ... */ }
    // "while" loop
    for condition { /* ... */ }
    // infinite loop
    for { /* ... */ }
    ```

    - The optional initialization statement is executed before the loop starts. If present, it must be **a simple statement**, that is, a short variable declaration, an increment or assignment, or a function call
    - The condition is a boolean expression that is evaluated at the beginning of each iteration of the loop
    - The post statement is executed after the body of the loop, then the condition is evaluated again
    - All of these 3 parts may be omitted
        - If there's no initialization and no post, the semicolons may be omitted
- `switch`
    - Cases are evaluated from top to bottom
    - `default` is optional
    - Cases do not fall through by default (use `fallthrough` override this)
    - A `switch` doesn't have to have an operand (tagless `switch`). Just list the cases, each of which is a boolean expression (equivalent to `switch true`)

    ```go
    switch coinflip() { // operand `coinflip()` is optional
        case "heads":
            head++ // don't fall through by default
        case "tails":
            tails++
            fallthrough
        default: // optional
            fmt.Println("land on edge!")
    }
    ```

- Labeled statement, `goto`