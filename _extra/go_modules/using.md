<!-- https://blog.golang.org/using-go-modules -->

- A module is a collection of Go packages stored in a file tree with a `go.mod` file at its root.
- The `go.mod` file defines the module’s *module path*, which is also the **import path** used for the root directory, and its dependency requirements, which are the other modules needed for a successful build.
# Creating a new module

```bash
# hello.go Outside GOPATH
# hello_test.go Outside GOPATH
go test
# PASS
# ok      _/Users/slackbuffer/Desktop/gomodules/hello     0.006s
```

- Because we are working outside `$GOPATH` and also outside any module, the `go` command knows no import path for the current directory and makes up a fake one based on the directory name: `_/Users/slackbuffer/Desktop/gomodules/hello`
- Make the current directory the root of a module by using `go mod init` and then try `go test` again

    ```bash
    go mod init example.com/hello
    go test
    # PASS
    # ok      example.com/hello       0.006s
    ```

    - The `go mod init` command wrote a `go.mod` file.
- The `go.mod` file only appears in the **root of the module**. Packages in subdirectories have import paths consisting of the module path plus the path to the subdirectory. 
    - For example, if we created a subdirectory `world`, we would not need to (nor want to) run `go mod init` there. The package would automatically be recognized as part of the `example.com/hello` module, with import path `example.com/hello/world`.
# Adding a dependency

```go
package hello

import "rsc.io/quote"

func Hello() string {
    return quote.Hello()
}
// run go test again
```

- The `go` command resolves imports by using the specific dependency module versions listed in go.mod. 
- When it encounters an import of a package not provided by any module in `go.mod`, the go command automatically looks up the module containing that package and adds it to `go.mod`, using the latest version. (“Latest” is defined as the latest tagged stable (non-prerelease) version, or else the latest tagged prerelease version, or else the latest untagged version.) 
    - `go test` resolved the new import `rsc.io/quote` to the module `rsc.io/quote v1.5.2`. 
    - It also downloaded two dependencies used by `rsc.io/quote`, namely `rsc.io/sampler` and `golang.org/x/text`. 
- Only direct dependencies are recorded in the `go.mod` file.
- A second `go test` command will not repeat this work, since the `go.mod` is now up-to-date and the downloaded modules are cached locally (in `$GOPATH/pkg/mod`).
- Adding one direct dependency often brings in other indirect dependencies too. The command `go list -m all` lists the current module and all its dependencies.

    ```bash
    go list -m all

    # example.com/hello
    # golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
    # rsc.io/quote v1.5.2
    # rsc.io/sampler v1.3.0
    ```

    - The current module, also known as the main module, is always the first line, followed by dependencies sorted by module path.
    - The `golang.org/x/text` version `v0.0.0-20170915032832-14c0d48ead0c` is an example of a [pseudo-version](https://golang.org/cmd/go/#hdr-Pseudo_versions), which is the go command's version syntax for a specific untagged commit.
- The `go` command maintains a file named `go.sum` containing the expected [cryptographic hashes](https://golang.org/cmd/go/#hdr-Module_downloading_and_verification) of the content of specific module versions:

    ```
    golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c h1:qgOY6WgZOaTkIIMiVjBQcw93ERBE4m30iBm00nkL0i8=
    golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
    rsc.io/quote v1.5.2 h1:w5fcysjrx7yqtD/aO+QwRjYZOKnaM9Uh2b40tElTs3Y=
    rsc.io/quote v1.5.2/go.mod h1:LzX7hefJvL54yjefDEDHNONDjII0t9xZLPXsUe+TKr0=
    rsc.io/sampler v1.3.0 h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=
    rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=
    ```

    - The go command uses the go.sum file to ensure that future downloads of these modules retrieve the same bits as the first download, to ensure the modules your project depends on do not change unexpectedly, whether for malicious, accidental, or other reasons. 
- Both `go.mod` and `go.sum` should be checked into **version control**.
- > https://github.com/golang/go/issues/24573#issuecomment-453578276
# Upgrading dependencies (minor version upgrades)
- With Go modules, versions are referenced with semantic version tags. 
- A semantic version has 3 parts: major, minor, and patch. For example, for `v0.1.2`, the major version is 0, the minor version is 1, and the patch version is 2.
- We can see we're using an untagged version of `golang.org/x/text`. Upgrade to the latest tagged version and test that everything still works:

    ```bash
    go get golang.org/x/text
    go test
    # PASS
    # ok      example.com/hello       0.007s
    ```

    - The `golang.org/x/text` package has been upgraded to the latest tagged version (v0.3.2). The `go.mod` file has been updated to specify v0.3.2 too. 
    - The `indirect` comment indicates a dependency is not used directly by this module, only indirectly by other module dependencies. 
    - See `go help modules` for details.
- Try upgrading the `rsc.io/sampler` minor version. The test failure shows that the latest version of `rsc.io/sampler` is incompatible with our usage.

    ```bash
    go get rsc.io/sampler
    # go: finding rsc.io/sampler v1.99.99
    # go: downloading rsc.io/sampler v1.99.99
    # go: extracting rsc.io/sampler v1.99.99
    go test
    # Failed

    go list -m -versions rsc.io/sampler
    # rsc.io/sampler v1.0.0 v1.2.0 v1.2.1 v1.3.0 v1.3.1 v1.99.99

    go get rsc.io/sampler@v1.3.1
    go test
    # OK
    cat go.mod
    ```

    - In general each argument passed to `go get` can take an explicit version; the default is `@latest`.
- > https://semver.org/
# Adding a dependency on a new major version
- `func Proverb` returns a Go concurrency proverb, by calling `quote.Concurrency`, which is provided by the module `rsc.io/quote/v3`. 
- Note that our module now depends on both `rsc.io/quote` and `rsc.io/quote/v3`:

    ```bash
    go list -m rsc.io/q...
    # rsc.io/quote v1.5.2
    # rsc.io/quote/v3 v3.1.0
    ```

- Each different major version (v1, v2, and so on) of a Go module uses a different module path: starting at v2, the path must end in the major version.
    - `v3` of `rsc.io/quote` is no longer `rsc.io/quote`: instead, it is identified by the module path `rsc.io/quote/v3`. 
        - This convention is called [semantic import versioning](https://research.swtch.com/vgo-import), and it gives incompatible packages (those with different major versions) different names.
    - In contrast, `v1.6.0` of `rsc.io/quote` should be backwards-compatible with `v1.5.2`, so it reuses the name `rsc.io/quote`.
        - In the previous section, `rsc.io/sampler v1.99.99` **should have** been backwards-compatible with `rsc.io/sampler v1.3.0`, but bugs or incorrect client assumptions about module behavior can both happen.
- The go command allows a build to include at most one version of any particular module path, meaning at most one of each major version: one `rsc.io/quote`, `one rsc.io/quote/v2`, one `rsc.io/quote/v3`, and so on. 
    - This gives module authors a clear rule about possible duplication of a single module path: it is impossible for a program to build with both `rsc.io/quote v1.5.2` and `rsc.io/quote v1.6.0`. 
- At the same time, allowing different major versions of a module (because they have different paths) gives module consumers the ability to **upgrade to a new major version incrementally**. 
    - In this example, we wanted to use `quote.Concurrency` from `rsc/quote/v3 v3.1.0` but are not yet ready to migrate our uses of `rsc.io/quote v1.5.2`. 
    - The ability to migrate incrementally is especially important in a large program or codebase.
# Upgrading a dependency to a new major version
- Complete our conversion from using `rsc.io/quote` to using only `rsc.io/quote/v3`.
- Reading the docs, we can see that Hello has become `HelloV3`:

    ```bash
    go doc rsc.io/quote/v3
    # package quote // import "rsc.io/quote/v3"

    # Package quote collects pithy sayings.

    # func Concurrency() string
    # func GlassV3() string
    # func GoV3() string
    # func HelloV3() string
    # func OptV3() string
    ```

- Update code and run test again.
# Removing unused dependencies
- We've removed all our uses of `rsc.io/quote`, but it still shows up in `go list -m all` and in our `go.mod` file.
- Building a single package, like with `go build` or `go test`, can easily tell when something is missing and needs to be added, but not when something can safely be removed. 
- Removing a dependency can only be done after checking all packages in a module, and all possible build tag combinations for those packages. An ordinary build command does not load this information, and so it cannot safely remove dependencies.
- The `go mod tidy` command cleans up these unused dependencies.
# Conclusion
- Go modules are the future of dependency management in Go. Module functionality is now available in all supported Go versions (that is, in Go 1.11 and Go 1.12).
- Workflow using Go modules:
    - `go mod init` creates a new module, initializing the `go.mod` file that describes it.
    - `go build`, `go test`, and other **package-building commands** add new dependencies to `go.mod` as needed.
    - `go list -m all` prints the current module’s dependencies.
    - `go get` changes the required version of a dependency (or adds a new dependency).
    - `go mod tidy` removes unused dependencies.
# `go help modules`
- Go 1.13 includes support for Go modules. Module-aware mode is active by default whenever a `go.mod` file is found in, or **in a parent of**, the current directory.
- For more fine-grained control, Go 1.13 continues to respect a temporary environment variable, `GO111MODULE`, which can be set to one of 3 string values: `off`, `on`, or `auto` (the default).
    - If `GO111MODULE=on`, then the `go` command requires the use of modules, never consulting `GOPATH`. 
        - We refer to this as the command being *module-aware* or running in "module-aware mode".
    - If `GO111MODULE=off`, then the `go` command never uses module support. Instead it looks in `vendor` directories and `GOPATH` to find dependencies.
        - We now refer to this as "`GOPATH` mode."
    - If `GO111MODULE=auto` or is unset, then the `go` command enables or disables module support based on the current directory. Module support is enabled only when the current directory contains a `go.mod` file or is **below** a directory containing a `go.mod` file.
- In module-aware mode, `GOPATH` no longer defines the meaning of imports during a build, but it still ***stores*** downloaded dependencies (in `GOPATH/pkg/mod`) and installed commands (in `GOPATH/bin`, unless `GOBIN` is set).
- A module is defined by a tree of Go source files with a `go.mod` file in the tree's **root directory**. The directory containing the `go.mod` file is called the *module root*.
- The module is the set of all Go packages in the module root and its subdirectories, but ***excluding*** subtrees with their own `go.mod` files.
- The *"module path"* is the import path prefix corresponding to the module root.
- The `go.mod` file can also specify replacements and excluded versions that only apply when building the module directly; they are ignored when the module is incorporated into a larger build.
- Once the `go.mod` file exists, no additional steps are required: go commands like 'go build', 'go test', or even 'go list' will automatically add new dependencies as needed to satisfy imports.
- By default, the `go` command satisfies dependencies by downloading modules from their sources and using those downloaded copies (after verification, as described in the previous section). 

- To allow interoperation with older versions of Go, or to ensure that all files used for a build are **stored together in a single file tree**, `go mod vendor` creates a directory named `vendor` in the root directory of the main module and stores there all the packages from dependency modules that are needed to support builds and tests of packages in the main module.
- To build using the main module's top-level `vendor` directory to satisfy dependencies (disabling use of the usual network sources and local caches), use `go build -mod=vendor`. 
    - Note that only the main module's **top-level** `vendor` directory is used; vendor directories in other locations are still ignored.
# `go help go.mod`
# `go help mod`