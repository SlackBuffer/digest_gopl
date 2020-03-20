# Array
- An array is a **fixed-length** sequence of zero or more elements of a **particular type**
- Individual array elements are accessed with the subscript notation
- `len` returns the number of elements in the array
- By default the elements of a new array variable are initially set to the zero value for the element type
- Use array literal to initialize an array with a list of values

    ```go
    var a [3]int // array of 3 integers
    var q [3]int = [3]int{1, 2, 3}
    var r [3]int = [3]int{1, 2}
    s := [...]int{1, 2, 3}

    type Currency int
    const (
        USD Currency = iota
        EUR
        GBP
        RMB
    )
    symbol := [...]string{RMB: "￥", USD: "$"}
    ```

    - If `...` appears in place of the length, the array length is determined by the number of initializers
    - Specify a list of index and value pairs and indices can appear in any order and some may be omitted
- The **size of an array is part of its type**
- The size must be a constant expression, that is, an expression whose value can be computed as the program is being **compiled**
- If an array's element is comparable then the array type is comparable too. `==` reports whether all corresponding elements are equal
- Arrays are **passed by value**, not by reference