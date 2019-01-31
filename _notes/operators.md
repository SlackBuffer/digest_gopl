# Binary operators
- Binary operators for arithmetic, logic, and comparison in order of decreasing precedence

    ```go
    *   /   %   <<  >>  &   &^
    +   -   |   ^
    ==  !=  <   <=  >   >=
    &&
    ||
    ```

- Operators at the same level **associate to the left**, so parentheses may be required for clarity
- Each operator in the **first two** lines of the table has a corresponding assignment operator like `+=`
- `+`, `-`, `*`, `/` may be applied to integer, floating-point, and complex numbers

    ```go
    fmt.Printf("%v %T\n", 5/1, 5/1) // 5 int
    fmt.Printf("%v %T", 5/1.0, 5/1.0) // 5 float64
    ```
    
- `%` applies only to integers
    - The behavior of `%` for negative numbers varies across programming languages
    - In Go, the sign of the reminder is always the same as the sign of the **dividend** (被取余数)
- If the result of an arithmetic operation, whether signed or unsigned, has more bits that can be represented in the result type, it's said to overflow. The high-order bits that do not fit are **silently discarded** (the sign may flip)
- `&&` and `||` have short-circuit behavior: if the answer is already determined by the value of the left operand, the right operand is not evaluated
## Comparability
- Two values of the same basic type (booleans, numbers, strings) may be compared using the `==` and `!=` operators
- Integers, floating-numbers, and strings are **ordered** by comparison operators
## Bitwise binary operators

```go
&   bitwise AND
|   bitwise OR
^   bitwise XOR (bitwise exclusive OR)
&^  bit clear (AND NOT)
<<  left shift
>>  right shift
```

- The first 4 treat their operands as bit patterns with no concept of arithmetic carry (进位) or sign
- `x<<n`, `x>>n`
    - The `n` determines the number of bit positions to shift and must be **unsigned**. The `x` operand may be unsigned or signed
    - Left shifts fill the vacated bits with zeros
    - Right shifts of unsigned numbers fill the vacated bits with zeros
    - Right shifts of signed numbers fill the vacated bits with copies of the **sign bit**
        - For this reason, it's important to use unsigned arithmetic when you're treating an integer as a bit pattern
    - > `x` and `n` don't have to of the same type
# Unary operators
- `+`: unary positive (no effect)
- `-`: unary negation
- For integers, `+x` is a shorthand for `0+x`; `-x` is a shorthand for `0-x`
- For floating-point and complex numbers, `+x` is just `x` and `-x` is the negation of `x`
- `^`: bitwise negation or complement
    - Return a value with each bit in its operand inverted
- `%c`, `%q` (with single quotes ), `%d` (numeric value) print runes