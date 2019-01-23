# Verbs
- https://golang.org/pkg/fmt/

    ```
    %d              decimal integer
    %x, %o, %b      integer in hexadecimal, octal, binary
    %f, %g, %e      floating-point number: 3.141593 3.141592653589793 3.141593e+00
    %t              boolean: true or false
    %c              rune (Unicode code point)
    %s              string
    %q              quoted string "abc" or rune 'c'
    %v              any value in a natural format
    %T              type of any value
    %%              literal percent sign (no operand)
    ```

- By convention, formatting functions whose names end in `f`, such as `log.Printf` and `fmt.Errorf`, use the formatting rules of `fmt.Printf`, whereas those whose names end in `ln` follow `Println`, formatting their arguments as if by `%v`, followed by a newline
- Variable type
    
    ```go
    fmt.Printf("%T", os.Args[:])
    ```