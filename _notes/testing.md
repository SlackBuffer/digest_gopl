# Testing
- A great deal of effort has been spent on techniques to make the complexity of programs manageable
- 2 techniques in particular stand out for their effectiveness
    1. Routine peer review of programs before they're deployed
    2. Testing
- Go's approach to testing relies on one command, `go test`, and a set of conventions for writing test functions that `go test` can run
    - The comparatively lightweight mechanism is effective for pure testing, and it extends naturally to benchmarks and systematic examples for documentation
- In practice, writing test code is not much different from writing the original programs itself
    - We write short functions that focus on one part of the task
    - We have to be careful of boundary conditions, think about data structures, and reason about what results a computation should produce from suitable inputs
- This is the same process as writing ordinary Go code; it needn't require new notations, conventions, and tools
## The `go test` tool
- The `go test` subcommand is a test driver for Go packages organized according to certain conventions
- In a package directory, files whose name end with `_test.go` are not part of the package ordinarily built by `go build` but are a part of it when built by `go test`
- Within `*_test.go` files, 3 kinds of functions are treated specially: tests, benchmarks, and examples
- A test function, which is a function whose name begins with `Test`, exercises some program logic for correct behavior
    - `go test` calls the test function and reports the result, which is either `PASS` or `FAIL`
- A benchmark function has a name beginning with `Benchmark` and measures the performance of some operation
    - `go test` reports the mean execution time of the operation
- An example function, whose name starts with `Example`, provides machine-checked documentation
- The `go test` tool scans the `*_test.go` files for these special functions, generates a temporary `main` package that calls them all in the proper way, builds and runs it, reports the results, and then clean up

- `go test -v exercises-the_go_programming_language/ch7/eval`
    - `-v` flag lets us see the printed output of the test, which is normally suppressed for a successful test
## `Test` functions
- Each test file must import the `testing` package
- Test functions have the following signature
	
    ```go
    func TestName(t *test.T) { /* ... */ }
    ```

- Test function names must begin with `Test`; the optional suffix `Name` must begin with a capital letter
- The `t` parameter provides methods for reporting test failures and logging additional information
- It's good practice to write the test first and observe that it triggers the same failure described by the bug report. Only then can we be confident that whatever fix we come up with addresses the right problem
    - As a bonus, running `go test` is usually quicker than manually going through the steps described in the bug report, allowing us to iterate more rapidly
- `-v` flag points the name and execution time of each test in the package
- If the test suite contains many slow tests, we may make even faster progress if we're selective about which ones we run. The `-run` flag, whose argument is a regular expression, causes `go test` to run only those tests whose function name matches the pattern (`go test -v -run="French|Canal"`)
- Once we've gotten the selected tests to pass, we should invoke `go test` with no flags to run the entire test suite one last time before we commit the change
- The style of table-driven testing in very common in Go
- The output of a failing test does not include the entire stack trace at the moment of the call to `t.Errorf`. Nor does `t.Errorf` cause a panic or stop the execution of the test
    - Tests are independent of each other. If an early entry in the table causes the test to fail, later table entries will still be checked, and thus we may learn about multiple failures during a single run
- When we really must stop a test function, perhaps because some initialization code failed or to prevent a failure already reported from causing a confusing cascade of others, use `t.Fatal` or `t.Fatalf`
    - These must be called from the same goroutine as the `Test` function, not from another one created during the test
- Test failure messages are usually of the form `"f(x) = y, want z`
    - When convenient, use actual Go syntax for the `f(x)` part
- Avoid boilerplate and redundant information
    - When testing a boolean function, omit the `want z` part since it adds no information
    - If `x`, `y`, or `z` is lengthy, print a concise summary of the relevant parts instead
### Randomized testing
- How to determine the output of a random input
    1. Write an alternative implementation of the function that uses a less efficient but simpler and clearer algorithm, and check that both implementations give the same result
    2. Create input values according to a pattern so that we know what output to expect
- Since randomized tests are nondeterministic, it's critical that the log of the failing test record sufficient information to reproduce the failure
    - For functions that accept complex inputs, it may be simpler toe **log the seed** of the pseudo-random number generator than to dump the entire input data structure. Armed with the seed value, we can easily modify the test to replay the failure deterministically
- By using the current time as a source of randomness, the test will explore novel inputs each time it's run
### Testing a command
- A package named `main` ordinarily produces an executable program, but it can be imported as a library too
    - Although a package name is `main` and it defines a `main` function, during testing this package acts as a library than exposes the function `TestEcho` to the test driver; its `main` function is ignored
- It's important that code being tested not call `log.Fatal` or `os.Exit`, since these will stop the process in its tracks
    - Calling these functions should be regarded as the exclusive right of `main`
- If something totally unexpected happens and a function panics, the test driver will recover, though the test will of course be considered a failure
- Expected errors such as those resulting from bad user input, missing files, or improper configuration should be reported by returning a non-nil `error` value
### White-box testing
- One way to categorizing tests is by the level of knowledge they require of the internal working of the package under test
- A black-box test assumes nothing about the package other than what is exposed by its API and specified by its documentation; the package's internals are opaque
    - Black-box tests are usually more robust, needing fewer updates as the software evolves
    - They also help the test author empathize with the client of the package and can reveal flaws in the API design
- A white-box test has privileged access to the internal functions and data structures of the packages and can make observations and changes that an ordinary client cannot
    - For example, a white-box test can check that the invariants of the package's data types are maintained after every operation
    - White-box tests can provide more detailed coverage of the trickier parts of the implementation
- This pattern (`storage2/storage_test.go`) can be used to temporarily save and restore all kinds of global variables, including command-line flags, debugging options, and performance parameters; to install and remove hooks that cause the production code to call some test code when something interesting happens; and to coax the production code into rare but important states, such as timeouts, errors, and even specific interleavings of concurrent activities
    - Using global variables in this way is safe only because `go test` does not normally run multiple tests concurrently
### External test packages
## Benchmark functions