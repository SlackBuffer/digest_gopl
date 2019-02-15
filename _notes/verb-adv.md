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

- `%08b`: `08` (adverb) modifies `%b` to pad the result with zeros to exactly 8 digits
- `[1]`, `#`

    ```go
    x := int64(0xdeadbeef)
    fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x) // 3735928559 deadbeef 0xdeadbeef 0XDEADBEEF
    ```

    - Explicit argument indexes
    - `#` emits prefix
- Floating-point
   1. `%g` prints floating-point values with the most compact representation that has adequate precision
   2. `%e` (exponent)
   3. `%f` (no exponent)
    - All 3 verbs allow field width and numeric precision to be controlled
- Space
    1. put spaces between bytes printing strings or slices in hex (`% x`, `% X`)
    <!-- 2. leave a space for elided sign in numbers (`% d`) -->
- `%U`
    - Unicode format
- `-`: 字段左对齐
- `width.precision`