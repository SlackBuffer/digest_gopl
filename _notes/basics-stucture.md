# Structure
- A Go program is stored in one or more files whose names end in `.go`. Each file begins with a `package` declaration that says what package the file is a part of
- The `import` declarations must follow the `package` declaration, followed by a sequence of package-level declarations of types, variables, constants, and functions, in any order
# Basics
- Go source files are always encoded in UTF-8 and Go text strings are conventionally interpreted as UTF-8
- Doesn't require **semicolons** at the ends of statements or declarations
    - Except where 2 or more appear on the same line
    - Newlines following certain tokens are converted into semicolons
        1. The opening brace `{` of the function must be on the same line as the end of the `func` declaration, not on a line by itself
        2. In the expression `x + y`, a newline is permitted after but not before the `+` operator
- All **indexing** in Go uses **half-open** intervals that include the first index but exclude the last
- `++`, `--`
    - **Prefix** only
    - **Statements**, not expressions (so `j = i++` is illegal)
- The **names** of Go function, variables, constants, **types**, statement labels, and packages follow a simple rule: a name begins with a letter (anything that Unicode deems a letter) or an underscore and may have any number of additional letters, digits, and underscores
    - No limit on name length
- Zero value
    - `0` for numbers, `false` for booleans, `""` for strings
    - `nil` for interfaces and **reference types** (slice, pointer, map, channel, function)
    - An aggregate type like an array or a struct has the zero value of all of its elements or fields
- 25 keywords

    ```go
    break       default     func    interface   select
    case        defer       go      map         struct    
    chan        else        goto    package     switch
    const       fallthrough if      range       type
    continue    for         import  return      var
    ```

- About 3 dozen predeclared names (like `int,` `true`) for built-in constants, types, and functions
    - These names are **not reserved**. Total fine to redeclaring

    ```go
    /* Constants */
    true false iota nil
    /* Types */
    int int8 int16 int32 int64
    uint uint8 uint16 uint32 uint64 uintptr
    float32 float64 complex128 complex64
    bool byte rune string error
    /* Functions */
    make len cap new append copy close delete
    complex real imag
    panic recover
    ```

- It's all bits at the bottom, but computers operate fundamentally on fixed-size number called **words**, which are interpreted as integers, floating-point numbers, bit sets, or memory address, then combine into larger aggregates
- Go's **type** fall into four categories
    1. basic types
       1. numbers
       2. strings
       3. boolean
    2. aggregate types
       1. arrays
       2. structs
    3. reference
       1. pointers
       2. slices
       3. maps
       4. functions
       5. channels
    4. interface types
- Reference types refer to program variables indirectly, so that the effect of an operation applied to one reference is observed by all copies of that reference
## `range`

```go
n := 0
for range "Hello, 世界" {
    n++
}
```

## Conversion
- `[]rune(str)` `[]byte(str)`