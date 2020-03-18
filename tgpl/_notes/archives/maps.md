- The hash table is an **unordered** collection of key/value paris in which all the **keys are distinct**, and the value associated with a given key can be retrieved, updated, or removed using **a constant number of key comparisons** on the average, no matter how large the hash table
- In Go, a map is **a reference to a hash table**, written `map[K]V`, where `K` and `V` are the types of its keys and values
    - All of the **keys** in a given map are **of the same type**. The key type `K` must be **comparable using `==`**, so that the map can test whether a given key is equal to one already within it
        - Through floating-point numbers are comparable, it's a bad idea to compare floats for equality (especially if `NaN` is a possible value)
    - All of the **values** are **of the same type**. There are no restrictions on the value type `V`
    - The keys need not be of the same type as the values
- Create maps

    ```go
    ages := make(map[string]int)
    ages["a"] = 1
    ages["b"] = 2

    // map literals
    ages := map[string]int{
        "a": 1,
        "b": 2,
    }

    // a new empty map
    m := map[string]int{}
    ```

- Map elements are accessed through the usual *subscript notation*

    ```go
    ages["a"] = 11
    delete(ages, "a")
    ```

    - The shorthand assignments forms `x += y` and `x++` also work for map elements
    - These operations are safe even the element isn't in the map. A map lookup using a key that isn't present returns the **zero value** for its type. So accessing a map element by subscripting always yields a value
- **A map elements is not a variable**, and we cannot take its address
    - One reason is that growing a map might cause rehashing of existing elements into new storage locations, thus potentially invalidating the address
- Use a `range`-based `for` loop to enumerate all the key/value paris in the map

    ```go
    for k, v := range ages {}
    ```

- The order of map iteration is unspecified, and different implementations might use a **different hash function**, leading to different ordering
- The **order** of map iteration is random, varying from one execution to another
    - This design is intentional, since it prevents programs from relying on any particular ordering where none is guaranteed
    - Making the sequence vary helps force programs to be robust across implementations
- To enumerate the key/value paris in order, we must sort the keys explicitly. Use the `Strings` function from the `sort` packages if the keys are string (this is a common pattern)

    ```go
    // var names []string
    names := make([]string, 0, len(ages))
    for name := range ages {
        names = append(names, name)
    }
    sort.Strings(names)
    for _, name := range names {
        fmt.Printf("%s\t%d\n", name, ages[name])
    }
    ```

- The zero value for a map type is `nil`, that is, a reference to no hash table at all

    ```go
    var ages map[string]int
    ages == nil // true
    len(ages) == 0 // true

    ages["sb"] = 28 // panic: assignment to entry in nil map
    ages = make(map[string]int) // allocate first
    ages["sb"] = 28 // ok
    ```

    - **Storing to a nil map cause a panic**. Must **allocate** the map before you can store into it
- Most operations on maps, including lookup, `delete`, `len`, and `range` loops, are safe to perform on a **nil map reference**, since it behaves like an empty map
- Distinguish between a nonexistent element and an element that happens to have the value zero

    ```go
    if age, ok := age["x"]; !ok { /* "x" is not a key in this map; age == 0 */ }
    ```

    - Subscripting a map in this context yields 2 values: the second is a boolean that reports whether the element was **present**
    - The boolean variable is often called `ok`, especially if it's immediately used in an `if` condition
- Maps cannot compared to each other; the **only legal comparison is with `nil`**. To test whether 2 maps contains the same keys and the same associated values, write a loop

    ```go
    func equal(x, y map[string]int) bool {
        if len(x) != len(y) {
            return false
        }
        for k, xv := range x {
            if yv, ok := y[k]; !ok || yv != xv {
                return false
            }
        }
        return true
    }
    ```

- Use a map whose keys are slices (not comparable)
   1. Define a helper function `k` that maps each key to a string, with the property that `k(x) == k(y)` if and only if we consider `x` and `y` are equivalent
   2. Create a map whose keys are strings, applying the helper function to each key before access the map

        ```go
        var m = make(map[string]int)

        func k(list []string) string { return fmt.Sprintf("%q", list) }
        func Add(list []string) { m[k(list)]++ }
        func Count(list []string) int { return m[k(list)] }
        ```

        - `k` use `fmt.Sprintf` to convert a slice of strings into a **single string** that is a suitable map key, quoting each slice element with `%q` to record string boundaries faithfully
            - The type of `k(x)` needn't be a string; any comparable type with the desired equivalence property will do, such as integers, arrays, or structs
        - The same approach can be used for any non-comparable key type. It's even useful for comparable key types when you want a **definition of equality other than `==`**, such as case-insensitive comparison for strings
- The value type of a map can itself be a composite type, such as a map or slice
- The idiomatic way to populate a map lazily is to initialize each value as its key appears for the first time
	
    ```go
    var graph = make(map[string]map[string]bool)
    func addEdge(from, to string) {
        edges := graph[from]
        if edges == nil {
            // populate lazily
            edges = make(map[string]bool)
            graph[from] = edges
        }
        edges[to] = true
    }
    func hasEdge(from, to string) bool {
        return graph[from][to]
    }
    ```
