- Go provides a mechanism called *reflection* to update variables and inspect their values at run time, to call their methods, and to apply the operations intrinsic to their representation, all without knowing their types at compile time.
- Reflection also lets us treat types themselves as first-class values.
- Go's reflection features increase the expressiveness of the language.
- Reflection is crucial to `fmt`, `encoding/json` (and alike), `text/template` (and alike).
    - Reflection is complex to reason about and not for casual use, so although these packages are implemented using reflection, they do not expose reflection in their own APIs.
# Why reflection
- Sometimes we need to write a function capable of dealing uniformly with values of types that don't satisfy a common interface, don’t have a known representation, or don’t exist at the time we design the function.
    - A familiar example is the formatting logic within `fmt.Fprintf`, which can usefully print an arbitrary value of any type, even a user-defined one.
- Implement a simple version called `Sprintf`: accept 1 argument; return the result as a string like `fmt.Sprintf`. Start with a type switch that tests whether the argument defines a `String` method, and call it if so. Then add switch cases that test the value’s dynamic type against each of the basic types—string, int, bool, and so on—and perform the appropriate formatting operation in each case.

    ```go
    func Sprint(x interface{}) string {
        type stringer interface {
            String() string
        }
        switch x := x.(type) {
            case stringer:
                return x.String()
            case string:
                return x
            case int:
                return strconv.Itoa(x)
            // ...similar cases for int16, uint32, and so on...
            case bool:
                if x {
                    return "true"
                }
                return "false"
            default:
                // array, chan, func, map, pointer, slice, struct
                return "???"
        }
    }
    ```

    - We could add more cases (`[]float64`, `map[string][]string`), but the number of types is infinite.
    - For named types, like `url.Values`, even if the type switch had a case for its underlying type `map[string][]string`, it wouldn’t match `url.Values` because the two types are not identical, and the type switch cannot include a case for each type like `url.Values` because that would require this library to depend upon its clients.
- Without a way to **inspect the representation of values of unknown types**, we quickly get stuck. What we need is reflection.
# `reflect.Type`, `reflect.Value`
- `reflect` package defines 2 important types, `Type` and `Value`.
- A `Type` represents a Go type. It is an ***interface*** with many methods for **discriminating** among types and **inspecting** their components, like the fields of a struct or the parameters of a function.
- The sole implementation of `reflect.Type` is the *type descriptor* (a set of **values that provide information about each type**, such as its name and methods; it's the same entity that identifies the dynamic type of an interface value).
- `reflect.TypeOf` accepts any `interface{}` and returns its **dynamic type** as a `reflect.Type`.

    ```go
    t := reflect.TypeOf(3)  // a reflect.Type
    fmt.Println(t.String()) // int

    // Type implements Stringer interface
	fmt.Println(t)			// int
    ``` 

    - The `TypeOf(3)` call above assigns the value `3` to the `interface{}` parameter. 
        - An assignment from a concrete value to an interface type performs an **implicit interface conversion**, which creates an **interface value** consisting of two components: its dynamic type is the operand’s type (`int`) and its dynamic value is the operand’s value (`3`).
    - Note that `reflect.Type` satisfies `fmt.Stringer`. Because printing the dynamic type of an interface value is useful for debugging and logging, `fmt.Printf` provides a **shorthand**, `%T`, that uses `reflect.TypeOf` internally.
    	
        ```go
        fmt.Printf("%T\n", 3)
        fmt.Printf("%s\n", reflect.TypeOf(3))
        ```
    
- Because `reflect.TypeOf` returns an interface value’s dynamic type, it always returns a **concrete type**.

    ```go
    var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // *os.File, not io.Writer
    ```

    - [ ] See later that `reflect.Type` is capable of representing interface types too
- A `reflect.Value` can hold a value of any type.
    - `reflect.Value` is a **struct**.
- The `reflect.ValueOf` function accepts any **`interface{}`** and returns a `reflect.Value` containing the interface’s **dynamic value**. As with `reflect.TypeOf`, the results of `reflect.ValueOf` are always **concrete**, but a `reflect.Value` **can hold interface values** too.

    ```go
    v := reflect.ValueOf(3)
    fmt.Println(v)			// 3
    
    fmt.Printf("%v\n", v)	// 3
    
    fmt.Println(v.String())                         // <int Value>
    fmt.Println(reflect.ValueOf("qwr").String())    // qwr

    t := v.Type()           // a reflect.Type
	fmt.Println(t.String()) // int
    ```

    - `reflect.Value` also satisfies `fmt.Stringer`, but unless the `Value` holds a **string** , the result of the `String` method reveals only the type. 
        - Instead, use the `fmt` package’s **`%v`** verb, which treats `reflect.Values` specially.
    - Calling the `Type` method on a `Value` returns its type as a `reflect.Type`.
- The **inverse** operation to `reflect.ValueOf` is the `reflect.Value.Interface` method. It returns an `interface{}` holding the same concrete value as the `reflect.Value`.

    ```go
    v := reflect.ValueOf(3) // a reflect.Value
    x := v.Interface()      // an interface{}
    i := x.(int)            // an int
    fmt.Printf("%d\n", i)   // 3
    ```

- A `reflect.Value` and an `interface{}` can both hold arbitrary values.
    - The difference is that an empty interface **hides** the representation and intrinsic operations of the value it holds and **exposes none of its methods**, so unless we know its dynamic type and use a type assertion to peer inside it, there is little we can do to the value within.
    - In contrast, a `Value` has many methods for **inspecting** its contents, regardless of its type.
- `ch12/format`
    - Instead of a type switch, we use `reflect.Value`'s `Kind` method to discriminate the cases. 
    - Although there are infinitely many types, there're only a finite number of kinds of **type**: 
        - the basic types Bool, String, and all the numbers; 
        - the aggregate types Array and Struct; 
        - the reference types Chan, Func, Ptr, Slice, and Map; 
        - Interface types;
        - `Invalid`, meaning no value at all. (The zero value of a `reflect.Value` has kind `Invalid`.)
    - `formatAtom` treats each value as an **indivisible** thing **with no internal structure**.
    - Since `Kind` is concerned only with the underlying representation, `format.Any` works for named types too.
# `Display`, a recursive value printer
- A debugging utility function `Display`: given an arbitrary complex value `x`, prints the complete structure of that value, labeling each element with the path by which it was found.
- Avoid exposing reflection in the API of a package where possible.
- Slices and arrays
    - `Len` returns the number of elements of a slice or array value.
    - `display` recursively invokes itself on each element on the sequence, appending the subscript notation "[i]" to the path.
- Structs
    - `NumField` reports the number of fields in the structs.
    - `Field(i)` returns the value of the `i`-th field as a `reflect.Value`.
    - To append the field selector notation ".f" to the path, we must obtain the `reflect.Type` of the struct and access the name of its `i`-th field. (`reflect.Type().Field(i).Name`)
    - Append the subscript notation "[key]" to the path.
- Maps
    - `MapKeys` returns a slice containing all the keys (of `reflect.Values` type) present in the map .
        - As usual when iterating over a map, the order is undefined.
    - `MapIndex(key)` returns the value corresponding to `key`.
        - Cutting a corner here. The type of a map key isn't restricted to the types `format.Any` handles best; arrays, structs, and interfaces can also be valid map keys.
- Pointers
    - `Elem` returns the variable pointed to by a pointer as a `reflect.Value`.
- Interfaces
    - Retrieve interface's dynamic value using `v.Elem()` and print its type and value.
- Although `reflect.Value` has many methods, only a few are safe to call on any given value.
- Even unexported fields are visible to reflection.
- `reflect.ValueOf` always returns a `Value` of a concrete type since it extracts the contents of an interface value.
- A `Value` obtained indirectly (`Value` in step 2 of `ch12/display`–`Display("&i", &i)`) may represent any value at all, including interfaces.