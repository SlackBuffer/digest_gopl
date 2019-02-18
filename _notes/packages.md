# Packages
- Standard library packages: https://golang.org/pkg
- Community packages: https://godoc.org
- The `golang.org/x/...` repositories hold packages designed and maintained by the Go team
    - The packages are not in the standard library because they're still under development or because they're rarely needed by the majority of Go programmers
- `go doc http.ListenAndServe`
- Packages in Go support modularity, encapsulation, separate compilation, and reuse
- The source code for a package resides in one or more `.go` files, usually in a directory whose name ends with the import path
    - The files of the `gopl.io/ch1/helloworld` package are stored in directory `$GOPATH/src/gopl.io/ch1/helloworld`
- Each package serves as a separate name space for its declarations
    - Within the `image` package, the identifier `Decode` refers to a different function than does the same identifier in the `unicode/utf16` package
    - To refer a function from outside its package, we must qualify the identifier to make explicit whether we mean `image.Decode` or `utf16.Decode`
- Packages let us hide information by controlling which names are visible outside the package, or exported
- Exported identifiers start with an upper-case letter
- Package names are always in lower case
## Imports
- Within a Go program, every package is identified by a unique string called **import path**
    - The language specification doesn't define where these strings come from or what they mean; it's up to the tools to interpret them
- When using the `go` tool, an import path denotes a directory containing one or more Go source file that together make up the package
- In addition to its import path, each package has a package name, which is short (and not necessarily unique) name that appears in its `package` declaration
- By convention, a package's name matches the last segment of its import path
    - An import declaration may specify an alternative name to avoid a conflict
- It's an error to import a package and then not refer to it
## Package initializing
- Package initialization begins by initializing package-level variables in the order in which they are declared, except that dependencies are resolved first
    - If the package has multiple `.go` files, they are initialized in the order in which the files are given to the compiler
    - The `go` tool sorts `.go` files by name before invoking the compiler
- One package is initialized at a time, in the order of imports in the program, dependencies first
    - So a package `p` importing `q` can be sure that `q` is fully initialized before `p`'s initialization begins
- Initialization proceeds from the bottom up; the `main` package is the last to be initialized 