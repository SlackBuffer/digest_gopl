- An **object** is simply a value or variable that has methods
- A method is a function associated with a (**any**) named type
- An object-oriented program is one that uses methods to express the properties and operations of each data structures so that **clients need not access the object's representation directly**

    ```go
    // Seconds method of type time.Duration
    const day = 24 * time.Hour
    day.Seconds() // 86400 
    ```

# Method declarations
- A method is declared with a variant of the ordinary function declaration in which **an extra parameter appears before the function name**
- The parameter attaches the function to the type of that parameter

    ```go
    func Distance(p, q Point) float64 {
        return math.Hypot(q.X-p.X, q.Y-p.Y)
    }
    func (p Point) Distance(q Point) float64 {
        return math.Hypot(q.X-p.X, q.Y-p.Y)
    }
	fmt.Println(p.Distance(q))
    ```

    - The extra **parameter** `p` is called the method's **receiver**, a legacy from early object-oriented languages that describes calling a method as "sending a message to an object"
- In Go, we don't use a special name like `this` or `self` for the receiver; we choose receiver names just as we would for any other parameter
    - Since the receiver name will be frequently used, it's a good idea to choose something short and to be consistent across methods
    - A common choice is **the first letter of the type name**
- In a method call, the receiver argument appears before the methods name. This parallels the declaration, in which the receiver **parameter** appear appears before the method name
- There's no conflict between the 2 declarations of functions called `Distance`
    - The first declares a package-level function `geometry.Distance`
    - The second declares a method of the type `Point`, so its name is `Point.Distance`
- The expression `p.Distance` is called a **selector**, because it selects the appropriate `Distance` method for the receiver `p` of type `Point`
- Selectors are also used to select fields of struct types, as in `p.X`
- Since methods and fields inhabit the same name space, declaring a method `X` on the struct type `Point` would be ambiguous and the compiler will reject it
- Since each type has its own name space for methods, we can use the name `Distance` for other methods so long as they belong to different types

    ```go
    type Path []Point
    // Distance returns the distance traveled along the path
    func (path Path) Distance() float64 {
        sum := 0.0
        for i := range path {
            if i > 0 {
                sum += path[i-1].Distance(path[i])
            }
        }
        return sum
    }
    ```

    - The complier determines which function to call based on both the method name and the type of the receiver
    - Go allows methods to be associated with any type
    - It's often convenient to define additional behaviors for simple types such as numbers, strings, slices, maps, and sometimes even functions
- Methods may be declared on any named type **defined in the same package**, so long as its underlying type is neither a **pointer** nor an **interface**
- All methods of a given type must have unique names, but different types can use the same name for a method; there's no need to qualify function names (for example, `PathDistance`) to disambiguate
- The first benefit to using methods over ordinary functions: methods names can be shorter
    - The benefit is magnified for calling originating outside the package, since they can use the shorter name and omit the package name
# Methods with a pointer receiver
- Because calling a function makes a copy of each argument value, if a function needs to **update a variable**, or if an argument is so **large** that we wish to avoid copying it, we must pass the address of the variable using a pointer
- The same goes for methods that need to update the receiver variable: we **attach them to the pointer type**

    ```go
    func (p *Point) ScaleBy(factor float64) {
        p.X *= factor
        p.Y *= factor
    }

    r := &Point{1, 2}
    r.ScaleBy(2)
    fmt.Println(*r) // "{2, 4}"

    p := Point{1, 2}
    (&p).ScaleBy(2)

    // shorthand
    p.ScaleBy(2)
    ```

    - The name of this method is `(*Point).ScaleBy`. Parentheses are necessary
- In a realistic program, **convention** dictates that if any method of `Point` has a pointer receiver, the **all methods** should have a pointer receiver, even ones that don't strictly need it
- **Named type and pointers to them** are the only types that may appear in a receiver declaration
- To avoid ambiguities, method declarations are not permitted on named types that are themselves pointer types

    ```go
    type P *int
    // receiver may be omitted
    func (P) f() { /* ... */ } // compile error: invalid receiver type
    ```

- If the `p` is a variable of type `Point` but the method requires a `*Point` receiver, we can use shorthand `p.ScaleBy(2)`. The compiler will perform an implicit `&p` on the variable
    - This works only for variables, including struct fields like `p.X` and array or slice elements like `perim[0]`
- We cannot call a `*Point` method method on a non-addressable `Point` receiver, because there's **no way to obtain the address of a temporary value**

    ```go
    Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal
    ```

- We can call a `Pointer` method with a `*Point` receiver, because there's a way to obtain the value from the address: just load the value pointed to by the receiver. The compiler inserts an implicit `*` operator for us

    ```go
    p := Point{1, 2}
    pptr := &p
    pptr.Distance(q)
    (*pptr).Distance(q)
    ```

- Summarize 3 cases (**compiler will perform implicit magic for *receiver***)
   1. The receiver argument has the same type as the receiver parameter (both have type `T` or both have type `*T`)
   2. The receiver argument is a variable of type `T` and the receiver parameter has type `*T`

        ```go
        p.ScaleBy(2) // implicit (&p)
        ```

        - The compiler implicitly takes the address of the variable
   3. The receiver has argument has type `*T` and the receiver parameter has type `T`

        ```go
        pptr.Distance(q) // implicit (*pptr)
        ```

        - The compiler implicitly ***dereferences*** the receiver, in other words, **loads the value**
- If all the methods of a named type `T` have a receiver type of `T` itself (not `*T`), it's safe to copy instances of that type; calling any of its methods necessarily makes a copy
    - For example, `time.Duration` values are liberally copied, including as arguments to functions
- If any method has a pointer receiver, you should avoid copying instances of `T` because doing so may violate internal invariants
    - For example, copying an instance of `bytes.Buffer` would cause the original and the copy to alias the same underlying array of bytes. Subsequent method calls would have unpredictable effects
## Nil is a valid receiver value
- Just as some functions allow nil pointers as arguments, so do some methods for their receiver, especially if `nil` is a meaningful zero value of the type, as with maps and slices

    ```go
    // an IntList is a linked list of integers
    // A nil *IntList represents the empty list
    type IntList struct {
        Value int
        Tail *IntList
    }
    // returns the sum of the list elements
    func(list *IntList) Sum() int {
        if list == nil {
            return 0
        }
        return list.Value + list.Tail.Sum()
    }
    ```

- When you define a type whose methods allow `nil` as a receiver value, it's worth pointing this out explicitly in its documentation comment
- In the final call to `Get` of urlvalue, the nil receiver behaves like an empty map
    - We could equivalent have written it as `Values(nil).Get("item")`, but `nil.Get("item")` would not compile because **the type of `nil`** has not been determined
    - Because `url.Values` is a map type and a map refers to its key/value pairs **indirectly**, any updates and deletions that `url.Values.Add` makes to the map elements are visible to the caller
    - However, as with ordinary functions, any changes a method makes to the reference itself, like setting it to `nil` or making it refer to a different map data structure, will not be reflected in the caller
# Composing types by struct embedding
- Struct embedding

    ```go
    type Point struct{ X, Y float64 }
    type ColoredPoint struct {
        Point
        Color color.RGBA
    }

    red := color.RGBA{255, 0, 0, 255}
    blue := color.RGBA{255, 0, 255, 255}
    var p = ColoredPoint{Point{1, 1}, red}
    var q = ColoredPoint{Point{5, 4}, red}

    fmt.Println(p.Distance(q.Point))
    p.ScaleBy(2)
    ```

    - We can call methods of the embedded `Point` field using a receiver of type `ColoredPoint`, even though `ColoredPoint` has no declared methods
    - The methods of `Point` have been **promoted** to `ColoredPoint`
- In this way, embedding allows complex types with many methods to be built up by the composition of several fields, each providing a few methods
- Coders familiar with class-based OO languages may be tempted to view `Point` as a base class and `ColoredPoint` as a subclass or derived class, or to interpret the relationship between these types as if a `ColoredPoint` is a `Point`. That would be a mistake
    - `Distance` has parameter of type `Point`, and `q` is not a `Point`, so although `q` does have an embedded field of that type, we must explicitly select it (`q.Point`). Attempting to pass `q` would be error (`p.Distance(q)`)
- A `ColoredPoint` is not a `Point`, but it "has a" `Point`, and it has 2 additional methods `Distance` and `ScaleBy` promoted from `Point`
    - If you prefer to think in terms of implementation, the embedded field instructs the complier to generate **additional wrapper methods** that delegate to the declared methods

        ```go
        func (p ColoredPoint) Distance(q Point) float64 {
            return p.Point.Distance(q)
        }
        func (p *ColoredPoint) ScaleBy(factor float64) {
            p.Point.ScaleBy(factor)
        }
        ```

- The type of an anonymous field may be a pointer to a named value, in which case fields and methods are **promoted indirectly** from the pointed-to object

    ```go
    type ColoredPoint struct {
        *Point
        Color color.RGBA
    }
    p := ColoredPoint{&Point{1, 1}, red}
    q := ColoredPoint{&Point{5, 4}, blue}
    fmt.Println(p.Distance(*q.Point))
    p.Point = p.Point // share the same Point
    ```

    - Adding another level of indirection lets us **share** common structures and vary the relationships between objects **dynamically**
- A struct may have more than one anonymous filed
- Methods can be declared only on named types (like `Point`) and pointers to them (`*Point`), but thanks to embedding, it's possible and sometimes useful for unnamed struct types to have methods too

    ```go
    var (
        mu sync.Mutex // guards mapping
        mapping = make(map[string]string)
    )
    func Lookup(key string) string {
        mu.Lock()
        v := mapping[key]
        mu.Unlock()
        return v
    }
    // functionally equivalent to
    var cache = struct {
        sync.Mutex
        mapping map[string]string
    } {
        mapping: make[string][string],
    }
    func Lookup(key string) string {
        cache.Lock()
        v := cace.mapping[key]
        cache.Unlock()
        return v
    }
    ```

    - The new variable gives more expressive names to the variables related to the cache. Promoting allows us to lock the `cache` with a self-explanatory syntax
# Methods values and expressions
- It's possible to separate the operations of selecting and calling a method
    - The selector `p.Distance` yields a method value, a function that binds a method (`Point.Distance`) to a specific value `p`
    - This function can then be invoked without a receiver value; it needs only the non-receiver arguments
- Method values are useful when a package's API calls for a function value, and the client's desired behavior for that function is to call a method on a specific receiver

    ```go
    type Rocket struct { /* ... */ }
    func (r *Rocket) Launch() { /* ... */ }
    r := new(Rocket)
    time.AfterFunc(10 * time.Second, func() { r.Launch() })
    time.AfterFunc(10 * time.Second, r.Launch )
    ```

- Method expression

    ```go
    type Point struct { X, Y float64 }
    func (p Point) Add(q Point) Point { return Point{p.X + 1.X, p.Y + q.Y }}
    func (p Point) Sub(q Point) Point { return Point{p.X - 1.X, p.Y - q.Y }}
    type Path []Point
    func (path Path) TranslateBy(offset Point, add bool) {
        var op func(p, q Point) Point
        if add {
            op = Point.Add
        } else {
            op = Point.Sub
        }
        for i := range path {
            path[i] = op(path[i], offset)
        }
    }
    ```

# Example: bit vector type
- A bit vector uses a slice of unsigned integer values or "words", each bit of which represents a possible element of the set
- The set contains `i` if `i`-th bit is set
- The `fmt` package treats types with a `String` method specially so that values of complicated types can display themselves in a user-friendly manner
    - Instead of printing the raw representation of the value, `fmt` calls the `String` method
        - > The mechanism relies on interfaces and type assertions
- Intset

    ```go
    var x IntSet
    // 1
    fmt.Println(&x) // {1 9 42 144}
    // 2
    fmt.Println(x.String()) // {1 9 42 144}
    // 3
    fmt.Println(x) // {[4398046511618 0 65536]}
    ```

    1. An `*IntSet` pointer does have a `String` method
    2. The compiler inserts the implicit `&` operation, giving us a pointer, which has the `String` method
    3. `IntSet` **value** does not have a `String` method; prints the representation of the struct instead
# Encapsulation
- A variable or method of an object is said to be encapsulated if it's inaccessible to clients of the object
- Encapsulation, sometimes called information hiding, is a key aspect of object-oriented programming
- Consequences of Go's name-based visibility control mechanism
   1. To encapsulate an object, we must use struct
        - Make the struct variable capitalized and its protected fields uncapitalized (avoid direct manipulation)
   2. The unit of encapsulation is the package, not the type
        - The fields of a struct type are visible to all code within the same package
- 3 benefits of encapsulation
   1. Clients cannot directly modify the object's variables, one need inspect fewer statements to understand the possible values of those variables
   2. Hiding implementation details prevents clients from depending on things that might change, which gives the designer greater freedom to evolve the implementation without breaking API compatibility
   3. Prevents clients from setting an object's variables arbitrarily
- Functions that merely access or modify internal values of a type, such as the methods of the `Logger` type from `log` package, are called getters and setters
- When naming a getter method, we usually omit the `Get` prefix
- This preference for brevity extends to all methods, not just field accessors, and to other redundant prefixes as well, such as `Fetch`, `Find`, and `Lookup`
- Go style does not forbid exported fields. Once exported, a field cannot be unexported without an incompatible change to the API
- Encapsulation is not always desirable
    - By revealing its representation as an `int64` number o nanoseconds, `time.Duration` lets us use all the usual arithmetic and comparison operations with durations, and even to define constants of this type