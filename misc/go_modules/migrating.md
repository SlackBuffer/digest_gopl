<!-- https://blog.golang.org/migrating-to-go-modules -->

# Without a dependency manager
# With a dependency manager
- For a Go project without a dependency management system, start by creating a `go.mod` file:

    ```bash
    git clone https://go.googlesource.com/blog
    cd blog
    go mod init golang.org/x/blog
    ```

    - Without a configuration file from a previous dependency manager, `go mod init` will create a `go.mod` file with only the module and `go` directives.
        - The `module` directive declares the module path, and the `go` directive declares the expected version of the Go language used to compile the code within the module.
    - We set the module path to `golang.org/x/blog` because that is its [custom import path](https://golang.org/cmd/go/#hdr-Remote_import_paths).
- Run `go mod tidy` to add the module's dependencies.
- `go mod tidy` added module requirements for all the packages transitively imported by packages in your module and built a `go.sum` with checksums for each library at a specific version.
- Finish by making sure the code still builds and tests still pass.

    ```bash
    go build ./...
    go test ./...
    ```

- Note that when `go mod tidy` adds a requirement, it adds the **latest** version of the module. If your `GOPATH` included an older version of a dependency that subsequently published a breaking change, you may see errors in `go mod tidy`, `go build`, or `go test`. If this happens, try downgrading to an older version with `go get` (for example, `go get github.com/broken/module@v1.1.0`), or take the time to make your module compatible with the latest version of each dependency.
# Tests in module mode
- Some tests may need tweaks after migrating to Go modules.
- If a test needs to write files in the package directory, it may fail when the package directory is in the module cache, which is read-only. In particular, this may cause `go test all` to fail. The test should copy files it needs to write to a temporary directory instead.
- If a test relies on relative paths (`../package-in-another-module`) to locate and read files in another package, it will fail if the package is in another module, which will be located in a versioned subdirectory of the module cache or a path specified in a `replace` directive. If this is the case, you may need to copy the test inputs into your module, or convert the test inputs from raw files to data embedded in `.go` source files.
- If a test expects `go` commands within the test to run in `GOPATH` mode, it may fail. If this is the case, you may need to add a `go.mod` file to the source tree to be tested, or set `GO111MODULE=off` explicitly.
# Imports and canonical module paths
- Each module declares its module path in its `go.mod` file. 
- Each import statement that refers to a package within the module must have the module path as a prefix of the package path. 
- However, the `go` command may encounter a repository containing the module through many different [remote import paths](https://golang.org/cmd/go/#hdr-Remote_import_paths). 
    - For example, both `golang.org/x/lint` and `github.com/golang/lint` **resolve to** repositories containing the code hosted at `go.googlesource.com/lint`. 
    - The `go.mod` file contained in that repository declares its path to be `golang.org/x/lint`, so only that path corresponds to a valid module.
- Go 1.4 provided a mechanism for declaring canonical import paths using [// import comments](https://golang.org/cmd/go/#hdr-Import_path_checking), but package authors did not always provide them. As a result, code written prior to modules may have used a non-canonical import path for a module without surfacing an error for the mismatch. 
- When using modules, the import path must match the canonical module path, so you may need to update import statements: for example, you may need to change `import "github.com/golang/lint"` to `import "golang.org/x/lint"`.
- A Go module with a major version above `1` must include a major-version suffix in its module path: for example, version `v2.0.0` must have the suffix `/v2`. 
    - However, `import` statements may have referred to the packages within the module without that suffix. For example, non-module users of `github.com/russross/blackfriday/v2` at `v2.0.1` may have imported it as `github.com/russross/blackfriday` instead, and will need to update the import path to include the `/v2` suffix.
# Conclusion
- Converting to Go modules should be a straightforward process for most users. Occasional issues may arise due to non-canonical import paths or breaking changes within a dependency. 