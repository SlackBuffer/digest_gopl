- `switch`
    - Cases are evaluated from top to bottom
    - `default` is optional
    - Cases do not fall through by default (use `fallthrough` override this)
    - A `switch` doesn't have to have an operand (tagless `switch`). Just list the cases, each of which is a boolean expression (equivalent to `switch true`)

    ```go
    switch coinflip() {
        case "heads":
            head++
        case "tails":
            tails++
        default:
            fmt.Println("land on edge!")
    }
    ```

- Labeled statement, `goto`