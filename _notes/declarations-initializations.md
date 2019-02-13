# Declaration
- `var name type = expression`
    - Either the type or the `=expression` part may be omitted, but not both
        1. If the type is omitted, it's determined by the initializer expression
        2. If the expression is omitted, the initial value is the zero value for the type
- Short variable declaration: `name := expression`
    - Only used within a function, to declare and initialize local variable
    - `:=` is a declaration; `=` is assignment
    - A short variable declaration does not necessarily declare all the variables on its left-hand side
        - If some of them were already declared in ***the same lexical block***, then, for those variables, the short variable declarations act like an **assignment** to those variables
            - Declarations in an outer block are ignored
    - A short variable declaration must declare **at least one new variable**
    - Short variable declarations with multiple initializer expressions should be used only when they help readability, such as for short and natural groupings like the initialization part of a `for` loop
- A `var` declaration tends to be reserved for local variables that need an explicit type that differs from that of the initializer expression

    ```go
    var boiling float64 = 100   // a float64
    ```

- Declaring a string variable

    ```go
    /* generally, use the first 2 forms */
    s := ""             // used only within function, not for package-level variables
    var s string        // relies on default initialization
    var s = ""          // rarely used except when declaring multiple variables
    var s string = ""   // explicit about the variable's type, necessary when it is not the same as that of the initial value
    ```

- Declare and optionally initialize a set of variables in a single declaration

    ```go
    var i, j, k int                 // int, int, int
    var b, f, s = true, 2.3, "four" // bool, float64, string
    ```

- `new(T)` creates an <mark>**unnamed variable**</mark> of type `T`, initializes it to the zero value of `T`, and returns its address, which is type `*T`
    - `new(T)` can be used in an expression
    - Each call to `new` returns a distinct unnamed variable with a different address
        - With **1 exception**: 2 variables whose type carries no information and is therefore of size zero, such as `struct{}` or `[0]int`, may, depending on the implementation, have the same address
    - `new` is a predeclared function, not a keyword, it's possible to redefine the name for something else
    - `new` is only a syntactic convenience. A variable created with `new` is no different from an ordinary local variable whose address is taken, except there's no need to invent (and declare) a dummy name

        ```go
        // same notion
        func newInt() *int {
            return new(int)
        }
        func newInt() *int {
            var dummy int
            return &dummy
        }
        ```

    - It's rarely used, because the most common unnamed variables are of struct, for which the struct literal syntax is more flexible
## Type declaration
- A `type` declaration defines a new **<mark>named type**</mark> that has the same underlying type as an existing type
    - `type name underlying-type`
- The named type provides a way to separate different and perhaps incompatible uses of the same underlying type so that they can't be mixed unintentionally

    ```go
    type Celsius float64
    type Fahrenheit float64
    ```

    - An explicit type conversion (`Celsius(t)`, `Fahrenheit(t)`) is required to convert from a `float64`. `Celsius(t)` and `Fahrenheit(t)` are conversions, not function calls. They don't change the value or representation in any way, but they make the change of meaning explicit
- Type declarations most often appear at package level. If the name is exported, it's accessible from other packages as well
- For every type `T`, there's a corresponding conversion operation `T(x)` that converts the value `x` to type `T`
- A conversion from one type to another is allowed if both have the same underlying type, or if both are unnamed pointer types that point to variables of the same underlying type
    - If `x` is assignable to `T`, a conversion is permitted but is usually redundant
- Conversions are also allowed between numeric types, and between string and some slice types
    - These conversions may change the representation of the value
        - Converting a floating-point number to an integer discards any fractional part
        - Converting a string to a `[]byte` slice allocates a copy of the string data
- The underlying type of a named type determines its structure and representation, and also the set of intrinsic operations it supports, which are the same as if the underlying type had been used directly
- Comparison operator like `==` and `<` can also be used to compare a value of a named type to another of the same type, or to a value of an unnamed type with the same underlying type
- A named type may provide **notational convenience** if it helps avoid writing out complex types (struct) over and over again
- Named types also make it possible to define new behaviors for values of the type. These behaviors are expressed as a set of functions associated with the type, called the type's methods

    ```go
    // Celsius **parameter** c appears before the function name
    func (c Celsius) String() string { return fmt.Sprintf("%g°C", c)}
    ```

    - Many types declare a `String` method of this form because it controls how values of the type appear when printed as a string by the `fmt` package
# Initialization
- Package-level variables are initialized **before `main` begins** (starts life with the value of its initializer expression, if any), and local variables are initialized as their declarations are encountered during function execution
- A set of variables can be initialized by calling a function that returns multiple values

    ```go
    var f, err = os.Open(name) // returns a file and an error
    f, err := os.Open(name)
    ```

## `init` function
- Declaration
    - `func init() { /* ... */}`
- Any file may contain any number of `init` functions
- Within each file, `init` functions are automatically executed when the program starts, in the order in which they are declared
- `init` function cannot be called or referenced, but otherwise they are normal functions
- Convenient to precompute a table of values
---
- `const`
    - The value of a constant must be a number, string, or boolean