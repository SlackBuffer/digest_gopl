# Scope
- The scope of a declaration is the part of the source code where a use of the declared name can refer to that declaration
- The scope of a declaration is a region of program text; it's a  compile-time property
- A *syntactic block* is a sequence of statements enclosed in braces like those that surround the body of a function or loop
- The block encloses its declarations and determines their scope. A name declared inside a syntactic block is not visible outside that block
- `Lexical block` is a generalized notion of blocks, and it includes other groupings of declarations that are not explicitly surrounded by braces
- There is a lexical block 
    - for the entire source code, called the **universe block**
    - for each package
    - for each file
    - for each `for`, `if`, and `switch` statement
    - for each case in a `switch` or `select` statement
    - for each explicit syntactic block
- A declaration's lexical block determines its scope
    - The declarations of built-in types, functions, and constants like `int`, `len`, and `true` are in the **universe block** and can be referred throughout the entire program
    - Declarations outside any function, that is, at **package level**, can be referred to from any file in the same package
    - Imported packages, such as `fmt`, are declared at the **file level**, so they can be referred to from the same file, but not from another file in the same package without another `import`
- The scope of a **control-flow label**, as used by `break`, `continue`, and `goto` statements, is the entire enclosing function
- A program may contain multiple declarations of the same name so long as each declaration is in a different lexical scope
- When the compiler encounters a reference to a name, it looks for a declaration, starting with the innermost enclosing lexical block and working up to the universe block
    - If the compiler finds no declaration, it reports an "undeclared name" error
    - If a name is declared in both an outer and an inner block, the inner declaration will be found first. In that case, the inner declaration is said to shadow or hide the outer one, making it inaccessible
- Implied lexical blocks
    - The `for` loop creates 2 lexical blocks: the explicit block for the loop body, and an implicit block that additionally encloses the variables declared by the initialization clause
        - The scope of a variable declared in the implicit block is the condition, post-statement (`i++`), and the body of the `for` statement

            ```go
            func main() {
                x := "hello"
                for _, x := range x {
                    x := x + 'A' - 'a'
                    fmt.Printf("%c", x) // "HELLO"
                }
            }
            ```

    - `if` statements and `switch` statements also create implicit blocks in addition to their body blocks

        ```go
        if x := f(); x == 0 {
            fmt.Println(x)
        } else if y := g(x); x ==y {
            fmt.Println(x, y)
        } else {
            fmt.Println(x, y)
        }
        ```

        - The second `if` statement is **nested** within the first, so variables declared within the first statement's initializer are visible within the second
    - There's a block for the condition and a block for each case body
- At the package level, the order in which declarations appear has no effect on their **scope**
- So a declaration may refer to itself (like functions, types!) or to another that follows it, letting us declare recursive or [ ] mutually recursive types and functions
    - The compiler will report an error if a **constant** or **variable** declaration refers to itself
# Visibility
- Function and other package-level entities may be declared in any order
- A entity whose name begins with an upper-case letter (**exported**)is visible and accessible outside its own package and may be referred to by other parts of the program
- An entity declared outside of a function is visible in all files of the package to which it belongs (package-level entity), as if the source code were all in a single file
- An entity declared within a function is local to that function
# Lifetime
- The lifetime of a variable is the range of time during execution when the variable can be referred to by other parts of the program; it's a run-time property
- The lifetime of a package-level variable is the entire execution of the program
- A new instance of a local variable is created each time the declaration is executed, and lives on util it becomes unreachable, at which point its storage may be recycled
    - Every package-level variable, and every local variable of each currently active function, can potentially be the start or root of a path to the variable in question, (variable) following pointers and other kinds of reference that ultimately lead to the variable
    - If no such path exists, the variable has become unreachable, so it can no longer affect the rest of the computation
- A compiler may choose to allocate local variables on the heap or on the stack. The choice is not determined by whether `var` or `new` was used to declare the variable

    ```go
    var global *int
    func f() {
        var x int
        x = 1
        global = &x
    }
    func g() {
        y := new(int)
        *y = 1
    }
    ```

    - `x` here must be heap-allocated because it's still reachable from `global` after `f` (function context exists on the stack) has returned, despite being declared as a local variable; we say `x` escapes from `f`
    - When `g` returns, the variable `*y` becomes unreachable and can be recycled. Since `*y` does not escape from `g`, it's safe for the compiler to allocate `*y` on the stack, even though it was allocated with `new`
- In closure cases, the lifetime of a variable is not determined by its scope