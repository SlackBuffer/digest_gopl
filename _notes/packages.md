# Packages
- Standard library packages: https://golang.org/pkg
- Community packages: https://godoc.org
- The `golang.org/x/...` repositories hold packages designed and maintained by the Go team
    - The packages are not in the standard library because they're still under development or because they're rarely needed by the majority of Go programmers
- `go doc http.ListenAndServe`
- Packages in Go support modularity, encapsulation, separate compilation, and reuse
- Within a Go program, every package is identified by a unique string called **import path**
    - Import paths are the strings that appear in `import` declarations
    - The language specification doesn't define where these strings come from or what they mean; it's up to the tools to interpret them
- When using the `go` tool, an import path denotes a directory containing one or more Go source file that together make up the package
- The source code for a package resides in one or more `.go` files, usually in a directory whose name ends with the import path
    - The files of the `gopl.io/ch1/helloworld` package are stored in directory `$GOPATH/src/gopl.io/ch1/helloworld`
- Each package serves as a separate name space for its declarations
    - Within the `image` package, the identifier `Decode` refers to a different function than does the same identifier in the `unicode/utf16` package
    - To refer a function from outside its package, we must qualify the identifier to make explicit whether we mean `image.Decode` or `utf16.Decode`
- Packages let us hide information by controlling which names are visible outside the package, or exported
    - Exported identifiers start with an upper-case letter
- When we change a file, we must recompile the file’s package and potentially all the packages that depend on it
- Go compilation is notably faster than most other compiled languages, even when built from scratch
    1. All imports must be explicitly listed at the beginning of each source file, so the compiler does not have to read and process an entire file to determines its dependencies
    2. The dependencies of a package form a directed acyclic graph
        - Because there are no cycles, packages can be compiled separately and perhaps parallel
    3. The object file for a compiled Go package records export information not just for the package itself, but its dependencies too
        - When compiling a package, the compiler must read one object file for each import but need not look beyond these files
## Naming
- **Package names are always in *lower case***
- When creating a package, keep its name short, but not so short as to be cryptic
- Be descriptive and unambiguous where possible
    - For example, don't name a utility package `util` when a name such as `imageutil` or `ioutil` is specific yet still concise
- Avoid choosing package names that are commonly used for related local variables, or you may compel the package's clients to use renaming imports
- Package names usually use the singular form
    - The standard packages `bytes`, `errors`, and `strings` use the plural to avoid hiding the corresponding types and, in the case of `go/types`, to avoid conflict with a keyword
- Avoid package names that already have other connotations
- When designing packages, consider how package name and member name of a qualified identifier work together, not hte member name alone
- *Single-type* packages such as `html/template` and `math/rand` expose one principal data type plus its methods, and often a `New` function to create instances
## Package declaration
- For packages you intend to share or publish, import paths should be globally unique
    - To avoid conflicts, the import paths of all packages other than those from the standard library should start with the internet domain name of the organization that owns or hosts the package; this also makes it possible to find packages
- A package declaration is required at the start of **every Go source file**
    - Its main purpose is to determine the default identifier for that package (called the package name) when it's imported by another package
- In addition to its import path, each package has a package name, which is short (and not necessarily unique) name that appears in its `package` declaration
- By convention, a package's name is the last segment of its import path
    - 2 packages may have the same name even though their import paths necessarily differ
- 3 major exceptions to the "last segment" convention
    1. A package **defining a command** (an executable Go program) always has the `main`, regardless of the package's import path
        - This is a signal to `go build` that it must invoke the linker to make an executable file
    2. Some files in the directory may have suffix `_test` on their package name if the file name ends with `_test.go`
        - Such a directory may define 2 package: the usual one, plus another one called an external test package
        - The `_test` suffix signals to `go test` that it must **build both packages**, and it indicates which file belong to each package
        - External test packages are used to avoid cycles in the import graph arising from dependencies of the test
    3. Some tools for dependency management append version number suffixes to package import path (`gopkg.in/yaml.v2`)
        - The package name excludes the suffix (`yaml`)
## Imports
- A Go source file may contain zero or more `import` declarations **immediately** after the `package` declaration and before the first non-import declaration
    - Each import declaration may specify the import path of a single package, or multiple packages in a parenthesized list
- (Parenthesized) Imported packages may be grouped by introducing blank lines; such groupings usually indicate different domains
    - The order is not significant, but by convention the lines of **each group** are sorted alphabetically (both `gofmt` and `goimports` will group and sort for you)
- An import declaration may specify an alternative name to avoid a conflict
	
    ```go
    import (
        "crypto/rand"
        mrand "math/rand"
    )
    ```

    - Called *renaming import*
    - The renaming import **affects only the importing file**. Other files (even ones in the same package) may import the package using its default name, or a different name
- A renaming import may be useful even when there's no conflict
    - If the name of the imported package is unwieldy, as is sometimes the case for automatically generated code, an abbreviated name may be more convenient
        - The same short name should be used consistently to avoid confusion
    - Choosing an alternative name can help avoid conflicts with common local variable names
        - e.g., in a file with many local variables named `path`, we might import the standard `"path"` packages as `pathpkg`
- Each import declaration establishes a dependency from the current package to the imported package
    - The `go build` tool reports an error if these dependencies form a cycle
- It's an error to import a package into a file but not refer to the name it defines within that file
- On occasion we must import a package merely for the **side effects** of doing so: evaluation of the initializer expressions and execution of its `init` functions
- Use a renaming import in which the alternative name is `_`, the blank identifier to suppress the "unused import" error
    - Known as *blank import*
    - As usual, **the blank identifier can never be referenced**
- The standard library provides decoders for GIF, PNG, and JPEG, and users may provide others, but to keep executables small, decoders are not included in an application unless explicitly requested
    - The `image.Decode` consults a table of supported formats. Each entry in the table specifies 4 things
        1. The name of the formats
        2. A string that is a prefix of all images encoded this way, used to detect the encoding
        3. A function `Decode` that decodes an encoded image
        4. `DecodeConfig` that decodes only the image metadata, such as its size and color space
    - An entry is added to the table by calling `image.ResisterFormat`, typically from within teh package initializer of the supporting package of each format
    	
        ```go
        func Decode(r io.Reader) (image.Image, error)
        func DecodeConfig(r io.Reader) (image.Config, error)
        func init() {
            const pngHeader = "\x89PNG\r\n\x1a\n"
            image.ResisterFormat("png", pngHeader, Decode, DecodeConfig)
        }
        ```
    
        - The effect is that an application need only blank-import the package for the format it needs to make the `image.Decode` function able to detect it
## Package initializing
- Package initialization begins by initializing package-level variables in the order in which they are declared, except that dependencies are resolved first
    - If the package has multiple `.go` files, they are initialized in the order in which the files are given to the compiler
    - The `go` tool sorts `.go` files by name before invoking the compiler
- One package is initialized at a time, in the order of imports in the program, dependencies first
    - So a package `p` importing `q` can be sure that `q` is fully initialized before `p`'s initialization begins
- Initialization proceeds from the bottom up; the `main` package is the last to be initialized 