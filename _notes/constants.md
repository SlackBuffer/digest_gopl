# Constants
- Constants are expressions whose value is known to the compiler and whose evaluation is guaranteed to occur at compile time, not at run time
- The underlying type of every constant is a **basic type**: boolean, string, or number
- `const`
- Prevents accidental or nefarious changes of value during program execution
- A sequence of constants can appear in one declaration; this would be appropriate for a group of related values
- Many computations on constants can be completely evaluated at compiled time, reducing the work necessary at run time and enabling other compiler optimizations
- Errors ordinarily detected at run time can be reported at compile time when their operands are constants
    - Such as integer division by zero, string indexing out of bounds, and any floating-point operation that would result in a non-finite value
- The result of all arithmetic, logical, and comparison operations applied to constant operands are themselves constants, as are the results of conversions and calls to certain built-in functions such as `len`, `cap`, `real`, `imag`, `complex`, and `unsafe.Sizeof`
- Since their values are known to the compiler, constant expressions may appear in types, specifically as the length of an array type
- A constant declaration may specify a type as well as a value, but in the absence of an explicit type, the type is inferred from the expression on the right-hand side
    - `time.Duration`, `time.Minute`
        - > Underlying type is `int64`
- When a sequence of constants is declared as a group, the right-hand side expression may be omitted for all but the first of the group, implying that the previous expression and its type should be used again

    ```go
    const (
        a = 1
        b
        c = 2 
        d
    )
    fmt.Println(a, b, c, d) // 1 1 2 2
    ```

## The constant generator `iota`
- A const declaration may use the constant generator `iota`, which is used to create a sequence of related values without spelling out each one explicitly
- In a `const` declaration, the value of `iota` begins at zero and increments by one for each item in the sequence
## Untyped constants
- Although a constant can have any of the basic data types like `int` or `float64`, including named basic types like `time.Duration`, many constants are not committed to a particular type
- The compiler represents these **uncommitted constants** with much greater numeric precision than values of basic types, and arithmetic on them is more precise than machine arithmetic; you may assume at least 256 bits of precision
- There are 6 flavors of these uncommitted constants, called untyped boolean, untyped integer, untyped rune, untyped floating-point, untyped complex, and untyped string
- By **deferring commitment**, untyped constants not only retain their higher precision until later, but they can participate in many more expressions than committed constants without requiring conversions
  
    ```go
    fmt.Println(YiB/ZiB) // 1024

    var x float32 = math.Pi
    var y float64 = math.Pi
    var z complex128 = math.Pi
    ```

    - ZiB and YiB are too big to store in any integer variable
    - `math.Pi` may be used wherever any floating-point or complex is needed. If `math.Pi` had been committed to specific type such as `float64`, the result would not be as precise, and type conversions would be required to use it when a `float32` or `complex128` is wanted
- For literals, syntax determines flavor
    - `0`, `0.0`, `0i`, and `\u0000` all denote constants of the same value but different flavors: untyped integer, untyped floating-point, untyped complex, and untyped rune
    - `true` and `false` are untyped booleans and string literals are untyped strings
- The choice of literal may affect the result of a constant division expression

    ```go
    var f float64 = 212
    fmt.Println((f - 32) * 5 / 9) // 100; (f - 32) * 5 is a float64
    fmt.Println(5 / 9 * (f - 32)) // 0; 5/9 is an untyped integer, 0
    fmt.Println(5.0 / 9.0 * (f - 32)) // 100; 5.0/9.0 is an untyped float
    ```

- Only constants can be untyped
- When an untyped constant is assigned to variable, or appears on the right-hand side of a variable declaration with an explicit type, the constant is implicitly converted to the type of that variable if possible
- Whether implicit or explicit, converting a constant from one type to another requires that the target type can represent the original value. Rounding is allowed for real and complex floating-point numbers

    ```go
    const deadbeef = 0xdeadbeef // untyped int with value 3735928559
    const b = float32(deadbeef) // float32 with value 37355928576 (rounded up)
    ```

- In a variable declaration without an explicit type (including short variable declarations), the favor of the untyped constant implicitly determines the default type of the variable

    ```go
    i := 0 // untyped integer; implicit int(0)
    r := '\000' // untyped rune; implicit rune('\000')
    f := 0.0 // untyped floating-point
    c := 0i // untyped complex
    ```

    - Untyped integers are converted to `int`, whose size is not guaranteed, but untyped floating-point and complex numbers are converted to the explicitly sized types `float64` and `complex128`
    - > The language has no unsized `float` and `complex` types analogous to unsized `int`, because it is very difficult to write correct numerical algorithms without knowing the size of one's floating-point data types
- To give the variable a different type, we must explicitly convert the untyped constant to the desired type or state the desired type in the variable declaration

    ```go
    var i = int8(0)
    var i int8 = 0
    ```

- These defaults are particularly important when converting an untyped constant to an **interface** value since they determine its dynamic type