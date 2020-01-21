# The `go` tool
- The Go toolchain converts a source program and the things it depends on into instructions in the **native machine language** of a computer, through a single `go` command with a number of subcommands
- The `go` tool combines the features of a diverse set of tools into one command set
    1. It's a package manager (analogous to `apt` or `rpm`) that answers queries about its inventory of packages, computes their dependencies, and downloads them from remote version-control systems
    2. It's a build system  that computes file dependencies and invokes compilers, assemblers, and linkers, although it's intentionally less complete than the standard Unix `make`
    3. It's a test driver
- To keep the need for configuration to a minimum, the `go` tool relies heavily on **conventions**
    - Given the name of a Go source file, the tool can find its enclosing package, because each directory contains a single package and the import path of a package corresponds to the directory hierarchy in the workspace
    - Given the import path of a package, the tool can find the corresponding directory in which it stores the object files
    <!-- - It can also find the URL of the server that hosts the source code repository -->
- `go help`
- `go help [command]`, `go help [topics]`
- `go run`
    - Compiles the source code from one or more source file whose name end in `.go`. ***Only compiles files specified in the arguments***
    - Links it with libraries
    - Runs the resulting executable file
- `go fmt` applies `gofmt` to all the files in the specified package, or the ones in the current directory by default
    - `gofmt` tool rewrites code into the standard format. Sorts the package names into alphabetical order
    - `goimports` manages the insertion and removal of import declarations as needed
- `go doc`
    - **`go doc http.ListenAndServe`**
    - `go doc html/template`
## Workspace organization
- The only configuration most user ever need is the `GOPATH` environment variable, which specifies the root of the workspace
- When switching to a different workspace, users update the value of `GOPATH`
	
    ```bash
    export GOPATH=$HOME/gobook
    go get gopl.io/...
    ```

- `GOPATH` has 3 subdirectories
    1. `src` holds source code
        - The path below `src` determines the import path or executable name
        - Each package resides in a directory whose name relative to `GOPATH/src` is the package's import path
        - A single `GOPATH` workspace may contain multiple version-control repositories beneath `src` (such as `gopl.io` or `golang.org`)
    2. `pkg` directory holds **installed package objects** (`.a` suffix)
        - As in the Go tree, each target operating system and architecture pair has its own subdirectory of `pkg` (`pkg/GOOS_GOARCH`)
        - > https://unix.stackexchange.com/questions/13192/what-is-the-difference-between-a-and-so-file
    3. `bin` holds compiled commands
        - Each command is named for its source directory, but only the final element, not the entire path
            - That is, the command with source in `DIR/src/foo/quux` is installed into `DIR/bin/quux`, not `DIR/bin/foo/quux`
            - The `"foo/"` prefix is stripped so that you can add `DIR/bin` to your `PATH` to get at the installed commands
        - If the `GOBIN` environment variable is set, commands are installed to the directory it names instead of `DIR/bin`
            - `GOBIN` must be an **absolute path**
- `GOROOT` specifies the root directory of the Go distribution, which provides all the p**ackages of the standard library**
    - The directory structure beneath `GOROOT` resembles that of `GOPATH`
        - The source files of the `fmt` package reside in the `$GOROOT/src/fmt` directory
    - Users never need to set `GOROOT` since, by default, the `go` tool will use the location where its was installed
- The `go env` command prints the effective values of the environment variables relevant to the toolchain, including the default values for the missing ones
    - `GOOS` specifies the target operating system
    - `GOARCH` specifies the target architecture
## Downloading packages
- When using the `go` tool, a package's import path indicates not only where to find it in the local workspace, but where to find it on the Internet so that `go get` can retrieve and update it
- `go get` can download a single package or an entire subtree or repository using the **`...`** notation
- The tool also computes and downloads all the dependencies of the initial packages
- Once `go get` has **downloaded** the packages, it **builds** them and then **installs** the libraries and commands
	
    ```bash
    go get github.com/golang/lint/golint

    $GOPATH/bin/golint gopl.io/ch2/popcount
    ```

- `go get` has support for popular code-hosting sites like GitHub, Bitbucket, and Lanuchpad and can make the appropriate requests to their version-control systems
    - For less well-known sites, you may have to indicate which version-control protocol to use in the import path, such as Git or Mercurial. Run `go help importpath` for the details
- The directories that `go get` creates are true clients of the remote repository, not just copies of the files, so you can **use version-control commands** to see a diff of local edits or to update to a different revision
	
    ```bash
    cd $GOPATH/src/golang.org/x/net
    git remote -v
    # origin  https://github.com/golang/net.git (fetch)
    # origin  https://github.com/golang/net.git (push)
    ```

    - The apparent domain name in the package's import path, `golang.org`, differs from the actual domain name of the Git server, `go.googlesrouce.com`. This is a feature of the `go` tool that lets packages use a **custom domain name** in their import path while being **hosted by a generic service** such as `googlesource.com` or `github.com`
    	
        ```bash
        go build gopl.io/ch1/fetch
        ./fetch https://golang.org/x/net/html | grep go-import
        # <meta name="go-import" content="golang.org/x/net git https://go.googlesource.com/net">

        curl http://gopl.io/ch1/helloworld?go-get=1
        ```

        - HTML pages beneath `https://golang.org/x/net/html` include the metadata shown below, which redirects the `go` tool to the Git repository at the actual site
        - HTTP requests from `go get` include the `go-get` parameter so that servers can distinguish them from ordinary browser requests
- If you specify the `-u` flag, `go get` will ensure that all packages it visits, including dependencies, are updated to their latest version before being built and installed. Without that flag, packages that already exists locally will not be updated
- `go get -u` generally retrieves the latest version of each package, which is convenient when you're getting started but may be inappropriate for deployed projects, where precise control of dependencies is critical for release hygiene
- The usual solution to this problem is to *vendor* the code, that is, to make a persistent local copy of all the necessary dependencies, and to update this copy carefully and deliberately
    - Prior to Go 1.5, this required changing those package's import paths, so the copy of `golang.org/x/net/html` would become `gopl.io/vendor/golang.org/x/net/html`
    - Most recent versions of the `go` tool support vendoring directly. See Vendor Directories in the output of the `go help gopath` command
## Building packages
- `go build` compiles each argument **package**
    - If the package is a **library**, the result is **discarded**; this merely checks that the package is free of compile errors
    - If the package is named `main`, `go build` invokes the linker to create an executable in the current directory; the **name of the executable** is taken from the last segment of the package's import path
- Since each directory contains one package, each executable, or command in Unix terminology, requires its own directory
    - These directories are sometimes children of a directory named `cmd`, such as the `golang.org/x/tools/cmd/godoc` command which serves Go package documentation through a web interface
- Packages may be specified by their import paths, or by a relative directory, which must start with a `.` or `..` segment even if this would not ordinarily be required. If no argument is provided, the **current directory is assumed**
- The following commands build the same package, though each writes the executable to the directory in which `go build` is run
	
    ```bash
    # current dir is assumed
    cd $GOPATH/src/gopl.io/ch1/helloworld
    go build

    # use import path
    cd anywhere
    go build gopl.io/ch1/helloworld

    # use relative path
    cd $GOPATH
    go build ./src/gopl.io/ch1/helloworld

    # not this one (wrong format of relative path)
    # Error: cannot find package "src/gopl.io/ch1/helloworld".
    cd $GOPATH
    go build src/gopl.io/ch1/helloworld
    ```

- Packages may also be specified as a list of file names, though this tends to be used only for small programs and one-off experiments
- If the package name is `main`, the executable name comes from the basename of the first `.go` file
- For throwaway programs, we want to run the executable as soon as we've built it. The `go run` command combines these 2 steps
	
    ```bash
    go run quoteargs.go one "two three" four\ five
    ```

    - The **first argument that doesn't end in `.go`** is assumed to be the beginning of the list arguments to the Go executable
- By default, the `go build` command builds the requested package and all its dependencies, then throws away all the compiled code except the final executable, if any
- The `go install` command is very similar to `go build`, except it **saves the compiled code for each package and command** instead of throwing it away
    - Compiled packages are saved beneath the `$GOPATH/pkg` directory corresponding to the `src` directory in which the source resides, and command executables are saved in the `$GOPATH/bin` directory (Many users put `$GOPATH/bin` on their executable search path). **Thereafter**, `go build` and `go install` do not run the compiler for those packages and commands if they have not changed, making subsequent builds much faster
    - For convenience, `go build -i` installs the packages that are dependencies of the build target
- Since compiled packages vary by platform and architecture, `go install` saves them beneath a subdirectory whose name incorporates the value of the `GOOS` and `GOARCH` environment variables
    - For example, on a Mac the `golang.org/x/net/html` package is compiled and installed in the file `golang.org/x/net/html.a` under `$GOPATH/pkg/darwin_amd64`
- It's straightforward to **cross-compile** a Go program, that is, to build an executable intended for a different system or CPU. Just set the `GOOS` or `GOARCH` variables during the build
	
    ```bash
    go build gopl.io/ch10/cross
    GOARCH=386 go build gopl.io/ch10/cross
    ```

- Some packages may need to compile different versions of the code for certain platforms or processors, to deal with low-level portability issues or to provide optimized versions of important routines, for instance
    - If a file name incudes an operating system or processor architecture name like `net_linux.go` or `asm_amd64.s`, then `go` tool will compile the file only when building for that target
    - Special comments called *build tags* give more fine-grained control
        - If a file contains this comment `+build linux darwin` before the package declaration (and its doc comment), `go build` will compile it only when building for Linux or Mac OS X. `+build ignore` says never to compile the file. For more details see the Build constraints section of the `go/build` package's documentation `go doc go/build`
## Documenting packages
- Go style strongly encourages good documentation of package APIs
- Each declaration of an exported package member and the package declaration itself should be immediately preceded by a comment explaining its purpose and usage
- Go doc comments are always **complete sentences**, and the first sentence is usually a summary that starts with the name being declared. Function parameters and other identifiers are mentioned **without quotation or markup**
	
    ```go
    // Fprintf formats according to a format specifier and writes to w.
    // It returns the number of bytes and any write error encountered.
    func Fprintf(w io.Writer, format string, a ...interface{})
    ```

    - The details of `Fprintf`'s formatting are explained in a doc comment associated with the `fmt` package itself
- A comment immediately preceding a `package` declaration is considered the doc comment for the package as a whole
    - There must be only one, though it may appear in any file
    - Longer package comments may warrant a file of their own. This file is usually called `doc.go`
- Go's convention favor **brevity** and **simplicity** in documentation as in all things. Many declarations can be explained in one well-worded sentence, and if the behavior is truly obvious, no comment is needed
- The `go doc` tool prints the declarations and doc comment of the entity specified on the command line, which may be a package, or a package member, or a method
	
    ```bash
    go doc time
    go doc time.Since
    go doc time.Duration.Seconds
    ```

    - The tool does not need complete import paths or correct identifier case
- `godoc` serves cross-linked HTML pages that provide the same information as `go doc` and much more
    - `godoc -http=:8080`, `godoc -http :8080`. Its `-analysis=type` and `-analysis=pointer` flags augment the documentation and the source code with the results of advanced static analysis
## Tools
### `go help gopath`
- If `DIR` is a directory listed in the `GOPATH`, a package with source in `DIR/src/foo/bar` can be imported as `"foo/bar"` and has its compiled form installed to `"DIR/pkg/GOOS_GOARCH/foo/bar.a"`
- Go searches each directory listed in `GOPATH` to find source code, but **new packages** are always downloaded into the **first directory** in the `GOPATH` list
- An example directory layout

    ```
    GOPATH=/home/user/go

    /home/user/go/
        src/
            foo/
                bar/               (go code in package bar)
                    x.go
                quux/              (go code in package main)
                    y.go
        bin/
            quux                   (installed command)
        pkg/
            linux_amd64/
                foo/
                    bar.a          (installed package object)
    ```

- When using modules, `GOPATH` is no longer used for resolving imports. However, it is still used to store downloaded source code (in `GOPATH/pkg/mod`) and compiled commands (in `GOPATH/bin`)
- Code in or below a directory named `"internal"` is importable only by code in the **directory tree rooted at the parent of `"internal"`**
	
    ```
    /home/user/go/
        src/
            crash/
                bang/              (go code in package bang)
                    b.go
            foo/                   (go code in package foo)
                f.go
                bar/               (go code in package bar)
                    x.go
                internal/
                    baz/           (go code in package baz)
                        z.go
                quux/              (go code in package main)
                    y.go
    ```

    - The code in `z.go` is imported as `"foo/internal/baz"`, but that import statement can only appear in source files in the subtree rooted at `foo`. The source files `foo/f.go`, `foo/bar/x.go`, and `foo/quux/y.go` can all import `"foo/internal/baz"`, but the source file `crash/bang/b.go` cannot
    - > https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit
- Go 1.6 includes support for using local copies of external dependencies to satisfy imports of those dependencies, often referred to as *vendoring*. Code below a directory named `"vendor"` is importable only by code in the **directory tree rooted at the parent of "vendor"**, and only using an import path that omits the prefix **up to and including** the `vendor` element
	
    ```
    /home/user/go/
        src/
            crash/
                bang/              (go code in package bang)
                    b.go
            foo/                   (go code in package foo)
                f.go
                bar/               (go code in package bar)
                    x.go
                vendor/
                    crash/
                        bang/      (go code in package bang)
                            b.go
                    baz/           (go code in package baz)
                        z.go
                quux/              (go code in package main)
                    y.go
    ```

    - The same visibility rules apply as for internal, but the code in `z.go` is imported as `"baz"`, not as `"foo/vendor/baz"`
- Code in vendor directories **deeper** in the source tree **shadows** code in higher directories
    - Within the subtree rooted at `foo`, an import of `"crash/bang"` resolves to `"foo/vendor/crash/bang"`, not the top-level `"crash/bang"`
- Code in vendor directories is **not subject to import path checking** (see `go help importpath`)
- When `go get` checks out or updates a git repository, it now also updates submodules
- Vendor directories do not affect the placement of new repositories being checked out for the first time by `go get`: those are always placed in the main `GOPATH`, never in a vendor subtree
    - When `go get` fetches a new dependency it never places it in the vendor directory. In general, moving code into or out of the vendor directory is the job of vendoring tools, not the `go` command
- > https://go.googlesource.com/proposal/+/master/design/25719-go15vendor.md
### `go help importpath`
- An import path (see `go help packages`) denotes a package stored in the local file system
- In general, an import path denotes either a standard package (such as `unicode/utf8`) or a package found in one of the **work spaces**
- An import path beginning with `./` or `../` is called a relative path. The toolchain supports relative import paths as a shortcut in 2 ways
    1. A relative path can be used as a shorthand on the command line
        - If you are working in the directory containing the code imported as `unicode` and want to run the tests for `unicode/utf8`, you can type `go test ./utf8` instead of needing to specify the full path
        - Similarly, in the reverse situation, `go test ..` will test `unicode` from the `unicode/utf8` directory
        - Relative patterns are also allowed, like `go test ./...` to test all subdirectories
    2. If you are compiling a Go program not in a work space, you can use a relative path in an import statement in that program to refer to nearby code also not in a work space
        - This makes it easy to experiment with small multipackage programs outside of the usual work spaces, but such programs cannot be installed with `go install` (there is no work space in which to install them), so they are rebuilt from scratch each time they are built
        - To avoid ambiguity, Go programs cannot use relative import paths within a work space
- Remote import paths
- For code hosted on other servers, import paths may either be qualified with the version control type, or the `go` tool can dynamically fetch the import path over https/http and discover where the code resides from a `<meta>` tag in the HTML
- To declare the code location, an import path of the form `repository.vcs/path`
- When a version control system supports multiple protocols, each is tried in turn when downloading
    - For example, a Git download tries `https://`, then `git+ssh://`
- By default, downloads are restricted to known secure protocols (e.g. `https`, `ssh`). To override this setting for Git downloads, the `GIT_ALLOW_PROTOCOL` environment variable can be set (For more details see `go help environment`)
- If the import path is not a known code hosting site and also lacks a version control qualifier, the `go` tool attempts to fetch the import over https/http and looks for a `<meta>` tag in the document's HTML `<head>`
- The `meta` tag has the form `<meta name="go-import" content="import-prefix vcs repo-root">`
    - The `import-prefix` is the import path corresponding to the repository root
        - It must be a prefix or an exact match of the package being fetched with `go get`
        - If it's not an exact match, another http request is made at the prefix to verify the `<meta>` tags match
    - The `vcs` is one of "bzr", "fossil", "git", "hg", "svn"
    - The `repo-root` is the root of the version control system containing a scheme and not containing a `.vcs` qualifier
- The `meta` tag should appear as early in the file as possible. In particular, it should appear before any raw JavaScript or CSS, to avoid confusing the go command's restricted parser
- A package statement is said to have an "import comment" if it is immediately followed (before the next newline) by a comment of one of these 2 forms

    ```go
	package math // import "path"
	package math /* import "path" */
    ```

- The `go` command will refuse to install a package with an import comment unless it is being referred to by that import path. In this way, import comments let package authors make sure the **custom import path is used** and **not a direct path to the underlying code hosting site**
- *Import path checking* is disabled for code found within vendor trees. This makes it possible to copy code into alternate locations in vendor trees without needing to update import comments
- Import path checking is also disabled when using modules. Import path comments are obsoleted by the `go.mod` file's module statement
- Go 1.4 Custom Import Path Checking: https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit
### `go help package`
### `go help gopath-get`
### `go help module-get`
### `go help goproxy`
### `go doc go/build`