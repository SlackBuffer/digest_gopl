# Structs
- A struct is an aggregate data type that groups together zero or more named values of arbitrary types as a single entity
- Each value is called a field
- All of these fields are collected into a single entity that can be copied as a unit, passed to functions and returned by them, stored in arrays, and so on
- Declare a struct type and a variable of that type

    ```go
    type Employee struct {
        ID int
        Name string
        Address string
        DoB time.Time
        Position string
        Salary int
        ManagerID int
    }
    var dilbert Employee
    ```

- The individual fields are accessed using dot notation
- A struct's fields are **variables**, so we may assign to a field or take its address and access it through a pointer

    ```go
    position := &dilbert.Position
    *position = "Senior " + *position
    ```

- ***The dot notation also works with a pointer to a struct***

    ```go
    var employeeOfTheMonth *Employee = &dilbert
    employeeOfTheMonth.Position += " (proactive team player)"
    // equivalent to
    (*employeeOfTheMonth).Position +=  " (proactive team player)"

    func EmployeeByID(id int) *Employee { /* ... */ }
    EmployeeById(dilbert.ID).Salary = 0
    ```

    - If the result type of `EmployeeByID` were changed to `Employee` instead of `*Employee`, the assignment statement would not compile since its left-hand side would not identify a variable (a struct literal)
- Fields are usually written one per line, with the field's name receding its type, but consecutive fields of the same type may be combined

    ```go
    type Employee struct {
        ID int
        Name, Address string
        DoB time.Time
        Position string
        Salary int
        ManagerID int
    }
    ```

    - Typically we only combine the declarations of related fields
- **Field order** is significant to type identity
    - Had we combined the declaration of the `Position` filed (also a string), or interchange `Name` and `Address`, we would be defining a different struct type
- A struct type may contain a mixture of exported and unexported fields
- Struct types usually appear within the declaration of a named type (like `Employee`)
- A named struct type can't declare a field of the same type `S`: an aggregate value cannot contain itself (An analogous restriction applies to arrays)
- But `S` may declare a fields of the pointer type `*S`, which lets us create recursive data structures like linked lists and trees
- The zero value for a struct is composed of the zero values of each of its fields
- It's usually desirable that the zero value be a natural or sensible default
    - e.g., in `bytes.Buffer`, the initial value of the struct is a ready-to-use empty buffer; the zero value of `sync.Mutex` is a ready-to-use unlocked mutex
    - Sometimes this sensible initial behavior happens for free, but sometimes the type designer has to work at it
- The struct type with no fields is called the empty struct, written `struct{}`
    - Its has size zero and carries no information but may be useful
    - Some Go programmers use it instead of `bool` as the value type of a map that represent a set, to emphasize that only the keys are significant, but the space saving is marginal and the syntax more cumbersome, so we generally avoid it

        ```go
        seen := make(map[string]struct{}) // set of strings
        if _, ok := seen[s]; !ok {
            seen[s] = struct{}{}
        }
        ```

## Struct literals
- 2 forms

    ```go
    // 1.
    type Point struct { X, Y int }
    p := Point{1, 2}

    // 2.
    anim := gif.GIF{LoopCount: nframes}
    ```

    1. The first form requires that a value be specified for every field, in the right order
        - Tends to be used only within the package that defines the struct type, or with smaller struct types for which there's an obvious filed ordering convention, like `image.Point{x, y}` or `color.RGBA{red, green, blue, alpha}`
    2. The second form lists some or all of the field names and their corresponding values
        - If a filed is omitted, it's set to the zero for its type
        - Order of fields doesn't matter
- The 2 forms cannot be mixed in the same literal
    - Nor can you use the order-based first form of literal to sneak around the rule that unexported identifies may not be referred to from another package
- Struct values can be passed as arguments to functions and returned from them
- For efficiency, larger struct types are usually passed to or returned from functions indirectly using a pointer
    - This is required if the function must modify its argument, since in a call-by-value language like Go, the called function receives only a copy of an argument, not a reference to the original argument
- Because structs are so commonly dealt with through pointers, it's possible to use this **shorthand notation** to create an initialize a `struct` variable and obtain its address

    ```go
    pp := &Point{1, 2}

    // equivalent to
    pp := new(Point) // returns an address
    *pp = Point{1, 2}
    ```

    - `&Point{1, 2}` can be used directly **within an expression**, such as a function call
## Comparing structs