<!-- https://github.com/golang/go/wiki/Modules -->

- An example of creating a module from scratch:
    1. Create a directory outside of your `GOPATH`, and optionally initialize VCS.
    2. Initialize a new module.
    3. Write code.
    4. Build and run.
- Daily workflow after that:
    - Add `import` statements to your `.go` code as needed.
    - Standard commands like `go build` or `go test` will automatically add new dependencies as needed to satisfy imports (updating `go.mod` and downloading the new dependencies).
    - When needed, more specific versions of dependencies can be chosen with commands such as `go get foo@v1.2.3`, `go get foo@master` (`foo@tip` with mercurial), `go get foo@e3702bed2`, or by **editing `go.mod` directly**.
- Common functionality:
    - `go list -m all` — View final versions that will be used in a build for all direct and indirect dependencies.
    - `go list -u -m all` — View available minor and patch upgrades for all direct and indirect dependencies.
    - `go get -u ./...` or `go get -u=patch ./...` (from module root directory) — Update all direct and indirect dependencies to latest minor or patch upgrades (pre-releases are ignored).
    - `go build ./...` or `go test ./...` (from module root directory) — Build or test all packages in the module.
    - `go mod tidy` — Prune any no-longer-needed dependencies from `go.mod` and **add** any dependencies needed for other combinations of OS, architecture, and build tags.
    - `replace` directive or `gohack` — Use a fork, local copy or exact version of a dependency.
    - `go mod vendor` — Optional step to create a `vendor` directory.
- Summarizing the relationship between repositories, modules, and packages:
    - A repository contains one or more Go modules.
    - Each module contains one or more Go packages.
    - Each package consists of one or more Go source files in a single directory.
- Modules must be semantically versioned according to semver, usually in the form `v(major).(minor).(patch)`.

- There are 4 directives: `module`, `require`, `replace`, `exclude`.
- A module declares its identity in its `go.mod` via the `module` directive, which provides the *module path*. 
    - The import paths for all packages in a module share the module path as a common prefix. 
    - The module path and the relative path **from the `go.mod` to a package's directory** together determine a package's import path.
- `exclude` and `replace` directives ***only*** operate on the current (“main”) module. `exclude` and `replace` directives in modules other than the main module are ignored when building the main module. 
- The `replace` and `exclude` statements, therefore, allow the main module complete control over its own build, without also being subject to complete control by dependencies. 
- The `replace` directive allows you to supply another import path that might be another module located in VCS (GitHub or elsewhere), or on your local filesystem with a relative or absolute file path. The new import path from the `replace` directive is used without needing to update the import paths in the actual source code.
- `replace` also can be used to inform the go tooling of the relative or absolute on-disk location of modules in a multi-module project, such as: `replace example.com/project/foo => ../foo`
    - If the right-hand side of a `replace` directive is a filesystem path, then the target must have a `go.mod` file at that location. If the `go.mod` file is not present, you can create one with `go mod init`.
- You can ***confirm*** you are getting your expected versions by running `go list -m all`, which shows you the actual final versions that will be used in your build including **taking into account `replace` statements**.

- In brief, to use vendoring with modules:
    - `go mod vendor` resets the main module's vendor directory to include all packages needed to build and test all of the module's packages based on the state of the `go.mod` files and Go source code.
    - By default, go commands like `go build` ignore the vendor directory when in module mode.
    - The `-mod=vendor` flag (e.g., `go build -mod=vendor`) instructs the go commands to use the main module's top-level vendor directory to satisfy dependencies. The go commands in this mode therefore ignore the dependency descriptions in `go.mod` and assume that the vendor directory holds the correct copies of dependencies. 
        - Note that ***only*** the main module's top-level vendor directory is used; vendor directories in other locations are still ignored.
    - Some people will want to routinely opt-in to vendoring by setting a `GOFLAGS=-mod=vendor` environment variable.
- Vendor **resets** the main module's vendor directory to include all packages needed to **build and test all the main module's packages**. 
    - It does not include test code for vendored packages.
    - `go mod vendor [-v]`: The `-v` flag causes vendor to print the names of vendored modules and packages to standard error.