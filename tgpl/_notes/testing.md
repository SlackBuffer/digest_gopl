# Testing
- "The realization came over me with full force that a good part of the remainder of my life was going to be spent in finding errors in my own programs."
- A great deal of effort has been spent on techniques to make the complexity of programs manageable. 2 techniques in particular stand out for their effectiveness
    1. Routine **peer review** of programs before they're deployed
    2. Testing
        - Testing, by which we implicitly mean *automated testing*, is the practice of writing small programs that check that the code under test (the production code) behaves as expected for certain inputs, which are usually either carefully chosen to exercise certain features or randomized to ensure broad coverage
- Go's approach to testing relies on one command, `go test`, and a set of conventions for writing test functions that `go test` can run
    - The comparatively lightweight mechanism is effective for pure testing, and it extends naturally to benchmarks and systematic examples for documentation
- In practice, writing test code is not much different from writing the original programs itself. We write short functions that focus on one part of the task. We have to be careful of boundary conditions, think about data structures, and reason about what results a computation should produce from suitable inputs. This is the same process as writing ordinary Go code; it **needn't** require new notations, conventions, and tools
    - Needn't acquire a whole new set of skills. This "low-tech" testing is just fine
## The `go test` tool
- The `go test` subcommand is a test driver for Go packages organized according to certain conventions
- In a package directory, files whose name end with `_test.go` are not part of the package ordinarily built by `go build` but are a part of it when built by `go test`
- Within `*_test.go` files, 3 kinds of functions are treated specially: tests, benchmarks, and examples
    - A *test function*, which is a function whose name begins with `Test`, exercises some program logic for correct behavior; `go test` calls the test function and reports the result, which is either `PASS` or `FAIL`
    - A *benchmark function* has a name beginning with `Benchmark` and measures the performance of some operation; `go test` reports the mean execution time of the operation
    - An *example function*, whose name starts with `Example`, provides machine-checked documentation
- The `go test` tool scans the `*_test.go` files for these special functions, generates a **temporary** `main` package that calls them all in the proper way, builds and runs it, reports the results, and then clean up
- `go test -v exercises-the_go_programming_language/ch7/eval`
    - `-v` flag lets us see the printed output of the test, which is normally suppressed for a successful test
    - `go test -v -run=TestFuncName .`
## `Test` functions
- Each test file must import the `testing` package. Test functions have the following signature
	
    ```go
    func TestName(t *testing.T) { /* ... */ }
    ```

- Test function names must begin with `Test`; the optional suffix `Name` must begin with a **capital letter**
- The `t` parameter provides methods for reporting test failures and logging additional information
- It's good practice to write the test first and observe that it triggers the same failure described by the bug report. Only then can we be confident that whatever fix we come up with addresses the right problem
    - As a bonus, running `go test` is usually quicker than manually going through the steps described in the bug report, allowing us to iterate more rapidly
- `-v` flag points the name and execution time of each test in the package
- If the test suite contains many slow tests, we may make even faster progress if we're selective about which ones we run. The `-run` flag, whose argument is a **regular expression**, causes `go test` to run only those tests whose function name matches the pattern (`go test -v -run="French|Canal"`)
    - Once we've gotten the selected tests to pass, we should invoke `go test` with no flags to run the entire test suite one last time before we commit the change
- The style of *table-driven* testing in very common in Go (`ch11/word2`)
- The output of a failing test does not include the entire stack trace at the moment of the call to `t.Errorf`. Nor does `t.Errorf` cause a panic or stop the execution of the test. Tests are independent of each other. If an early entry in the table causes the test to fail, later table entries will still be checked, and thus we may learn about multiple failures during a single run
- When we really must stop a test function, perhaps because some initialization code failed or to prevent a failure already reported from causing a confusing cascade of others, use `t.Fatal` or `t.Fatalf`. These must be called from the **same goroutine** as the `Test` function, not from another one created during the test
- Test failure messages are usually of the form `"f(x) = y, want z"`, where `f(x)` explains the attempted operation and its input, `y` is the actual result, and `z` the expected result
    - When convenient, use actual Go syntax for the `f(x)` part
    - Displaying `x` is particularly important in a table-driven test, since a given assertion is executed many times with different values
- Avoid boilerplate and redundant information
    - When testing a boolean function, omit the `want z` part since it adds no information. If `x`, `y`, or `z` is lengthy, print a concise summary of the relevant parts instead
    - The author of a test should strive to help the programmer who must diagnose a test failure
### Randomized testing
- Table-driven tests are convenient for checking that a function works on inputs carefully selected to exercise interesting cases in the logic. *Randomized testing*, explores a broader range of input by constructing inputs at random
- How to determine the output of a random input
    1. Write an alternative implementation of the function that uses a less efficient but simpler and clearer algorithm, and check that both implementations give the same result
    2. Create input values according to a pattern so that we know what output to expect (`ch11/rand_test`)
- Since randomized tests are nondeterministic, it's critical that the log of the failing test record sufficient information to reproduce the failure
    - For functions that accept complex inputs, it may be simpler to **log the seed** of the pseudo-random number generator than to dump the entire input data structure. Armed with the seed value, we can easily modify the test to replay the failure deterministically
    - By using the current time as a source of randomness, the test will explore novel inputs each time it's run, over the entire course of its lifetime. This is especially valuable if your project uses an automated system to run all its tests periodically
### Testing a command
- A package named `main` ordinarily produces an executable program, but it can be imported as a library too
    - Although a package name is `main` and it defines a `main` function, during testing this package acts as a library than exposes the function `TestEcho` to the test driver; its `main` function is ignored
- It's important that code being tested not call `log.Fatal` or `os.Exit`, since these will stop the process in its tracks
    - Calling these functions should be regarded as the exclusive right of `main`
- If something totally unexpected happens and a function panics, the test driver will recover, though the test will of course be considered a failure
- Expected errors such as those resulting from bad user input, missing files, or improper configuration should be reported by returning a non-nil `error` value
### White-box testing
- One way to categorizing tests is by the level of knowledge they require of the internal working of the package under test
- A *black-box* test assumes nothing about the package other than what is exposed by its API and specified by its documentation; the package's internals are opaque
    - Black-box tests are usually more robust, needing fewer updates as the software evolves
    - They also help the test author empathize with the client of the package and can reveal flaws in the API design
    - `TestIsPalindrome` calls only the exported function and is thus a black-box test
- A *white-box* test has privileged access to the internal functions and data structures of the packages and can make observations and changes that an ordinary client cannot
    - For example, a white-box test can check that the invariants of the package's data types are maintained after every operation
    - White-box tests can provide more detailed coverage of the trickier parts of the implementation
    - `TestEcho` calls `echo` and updates the global variable `out`, both of which are unexported, making it a white-box test
- While developing `TestEcho`, we modified the `echo` function to use the package level variable `out` when writing its output, so that the test could replace the standard output with an alternative implementation that records the data for later inspection. Using the same technique, we can replace other parts of the production code with easy-to-test **"fake"** implementations      
    - The advantage of fake implementations is that they can be simpler to configure, more predictable, more reliable, and easier to observe. They also avoid undesirable side effects such as updating a production database or changing a credit card
    - `ch11/storage1`
- `ch11/storage2`
    - This pattern can be used to temporarily save and restore all kinds of global variables, including command-line flags, debugging options, and performance parameters; to install and remove hooks that cause the production code to call some test code when something interesting happens; and to coax the production code into rare but important states, such as timeouts, errors, and even specific interleavings of concurrent activities
    - Using global variables in this way is safe only because `go test` does not normally run multiple tests concurrently
### External test packages
- Consider the package `net/url`, which provides a URL parser, and `net/http`, which provides a web server and HTTP client library. The higher-level `net/http` depends on the lower-level `net/url`. One of the tests in `net/url` is an example demonstrating the interaction between URLs and HTTP client library. In other words, a test of the lower-level package imports the higher-level package
  - Declaring this test function in the `net/url` package would create a cycle in the package import path. Go specification **forbids** import cycles
    ![](src/cycle.jpg)
- We resolve the problem by declaring the test function in an *external test package*, that is, in a file in the `net/url` directory whose package declaration reads `package url_test` (would be `package url` normally)
    - The extra suffix `_test` is a signal to `go test` that it should build an additional package containing just these files and run its tests. It may be helpful to think of this external test package as if it had the import path `net/url_test`, but it cannot be imported under this or any other name
    ![](src/break_cycle.jpg)
    - By avoiding import cycles, external test packages allow tests, especially integration tests (which test the interaction of several components), to import other packages freely, exactly as an application would
- We can use the `go list` tool to summarize which Go source files in a package directory are production code, in-package tests, and external tests
	
    ```bash
    # use fmt package as an example
    go list -f={{.GoFiles}} fmt
    go list -f={{.TestGoFiles}} fmt # [export_test.go]
    go list -f={{.XTestGoFiles}} fmt
    ```

    - `GoFiles` is the list of files that contain the production code; these are the files that `go build` will include in the application
    - `TestGoFiles` is the list of files that also belong to the `fmt` package, but these files, whose names all end in `_test.go`, are included only when building tests
    - `XTestGoFiles` is the list of files that constitute the external test package, `fmt_test`, so these files must import the `fmt` package on order to use it. They are only included during testing
- Sometimes an external test package may need privileged access to the internals of the package under test, if for example a white-box test must live in a separate package to avoid an import cycle. In such cases, we use a trick: we add declarations to an in-package `_test.go` file to expose the necessary internals to the external test. This file thus offers the test a back door to the package
    - If the source file exists only for this purpose and contains no tests itself, it's often called `export_test.go`
- The implementation of the `fmt` package needs the functionality of `unicode.IsSpace` as part of `fmt.Scanf`. To avoid creating an undesirable dependency, `fmt` does not import the `unicode` package and its large tables of data; instead, it contains a simpler implementation, which it calls `isSpace`
- To ensure that the behavior of `fmt.isSpace` and `unicode.IsSpace` do not drift apart, `fmt` contains a test. It's an external test, and thus it cannot access `isSpace` directly, so `fmt` opens a back door to it by declaring an exported variable that holds the internal `isSpace` function. This is the entirety of the `fmt` package's `export_test.go` file
	
    ```go
    // /usr/local/go/src/fmt/export_test.go
    package fmt
    var IsSpace = isSpace
    ```

    - This trick can also be used whenever an external test need to use some of the techniques of white-box testing
### Writing effective tests
- Many newcomers to Go are surprised by the minimalism of Go's testing framework. Other languages' frameworks provide mechanisms for identifying test functions (often using reflection or metadata), hooks for performing "setup" and "teardown" operations before and after the tests fun, and libraries of utility functions for asserting common predicates, comparing values, formatting error messages, and aborting a failed test (often using exceptions)
    - Although these mechanisms can make tests very concise, the resulting tests often seem like they are in a foreign language
    - Furthermore, although they may report PASS or FAIL correctly, their manner may be unfriendly to the unfortunate maintainer, with cryptic failure messages like "assert: 0 == 1" or page after age of stack traces
- Go expects test authors to do most of this work themselves, defining functions to avoid repetition, just as they would for ordinary programs. The process of testing is not one of the role form filling; a test has a user interface too, albeit one whose only users are also its maintainers
    - A good test does not explode on failure but prints a clear and succinct description of the symptom of the problem, and perhaps other relevant facts about the context
    - Ideally, the maintainers should not need to read the source code to decipher a test failure
    - A good test should not give up after one failure but should try to report several errors in a single run, since the pattern of failures may itself be revealing
- The following assertion function suffers from *premature abstraction*

    ```go
    func assertEqual(x, y int) {
        if x != y {
            panic(fmt.Sprintf("%d != %d", x, y))
        }
    }
    func TestSplit(t *testing.T) {
        words := strings.Split("a:b:c", ":")
        assertEqual(len(words), 3)
        // ...
    }
    ```

    - By treating the failure of this particular test as a mere difference of 2 integers, we forfeit the opportunity to provide meaningful context
    - Provide a better message by starting from the concrete details. Only once repetitive patterns emerge in a given test suite is it time to introduce abstractions

        ```go
        func TestSplit(t *testing.T) {
            s, sep := "a:b:c", ":"
            words := strings.Split(s, sep)
            if got, want := len(words), 3 got != want {
                t.Errorf("Split(%q, %q) returned % words, want %d", s, sep, got, want)
            }
        }
        ```

        - Once we've written a test like this, the natural next step is often not to define a function to replace the entire `if` statement, but to execute the test in a loop in which `s`, `sep`, and `want` vary
- The key to a good test is to start by implementing the concrete behavior that you want and only then use functions to simplify the code and eliminate repetition
### Avoid brittle tests
- A test that spuriously fails when a sound change was made to the program is called *brittle*
    - The most brittle tests, which fail or for almost any change to the production code, good or bad, are sometimes called *change detector* or *status quo* tests
- The easiest way to avoid brittle test is to check only the properties you care about. Test the program's simpler and more stable interfaces in preference to its internal functions. Be selective in the assertions
    - Don't check for exact string matches, but look for relevant substrings that will remain unchanged as the program evolves
- It's often worth writing a substantial function to **distill** a complex output down to its essence so that assertions will be reliable
## Coverage
- > Edsger Dijkstra: "Testing shows the presence, not the absence of bugs."
- The degree to which a test suite exercises the package under test is called the test's *coverage*
- Coverage cannot be quantified directly–the dynamics of all but the most trivial programs are beyond precise measurement–but there're heuristics that help us direct the testing effort to where they are more likely to be useful
- *Statement coverage* is the simplest and most widely used of these heuristics. The statement coverage of a test suite is the fraction of source statements that are executed at least once during the test
- `ch7/eval`
	
    ```bash
    go test -v -run=Coverage digest_gopl/ch7/eval

    go test -run=Coverage -coverprofile=b.out digest_gopl/ch7/eval

    go test -run=Coverage -coverprofile=c.out -covermode=count digest_gopl/ch7/eval
    # process the log and generate an HTML report
    go tool cover -html=c.out
    ```

    - Go's `cover` tool is integrated into `go test`. The `go tool` command runs one of the executables from the Go toolchain. These programs live in the directory `$GOROOT/pkg/tool/${GOOS}_${GOARCH}`. Thanks to `go build`, we rarely need to invoke them directly
        - `go tool cover` displays the usage message of the coverage tool
    - `-coverprofile` enables the collection of coverage data by *instrumenting* the production code. This is, it modifies a copy of the source so that before each block of statements is executed, a boolean variable is set, with one variable per block. Just before the modified program exits, it writes the value of each variable to the specified file `c.out` and prints a summary of the fraction of statement that were executed
        - If all you need is the summary, use `go test -cover`
    - If `go test` is run with the `-covermode=count` flag, the instrumentation for each block increments a counter instead of setting a boolean. The resulting log of execution counts of each block enables quantitative comparisons between hotter blocks and colder ones
- Achieving 100% statement coverage sounds like a noble goal, but it's not usually feasible in practice, nor is it likely to be a good use of effort
    - Just because a statement is executed does not mean it's bug-free; statements containing complex expressions must be executed many times with different inputs to cover the interesting cases
    - Some statements, like the `panic` statements above, can never be reached. Others, such as those that handle esoteric errors, are hard to exercise but rarely reached in practice
- Testing is fundamentally a pragmatic endeavor, a trade-off between the cost of writing tests and the cost of failures that could have been prevented by tests. Coverage tools can help identify the weakest spots, but devising good test cases demands the same rigorous thinking as programming in general
## Benchmark functions
- Benchmarking is the practice of measuring the performance of a program on a fixed workload
- In Go, a benchmark function look like a test function, but with the `Benchmark` prefix and a `*testing.B` parameter that provides most of the same methods as a `*testing.T`, plus a few extra related to performance measurement
    - It also exposes an integer field `N`, that specifies the number of times to perform the operation being measured
    	
        ```go
        // ch11/word2/word_test.go
        func BenchmarkIsPalindrome(b *testing.B) {
            for i := 0; i < b.N; i++ {
                IsPalindrome("A man, a plan, a canal: Panama")
            }
        }
        // go test -bench=.
        // go test -bench=IsPalindrome

        // BenchmarkIsPalindrome-8         10000000               151 ns/op
        // PASS
        // ok      digest_gopl/ch11/word2  1.679s
        ```
    
    - Unlike tests, by default, no benchmarks are run. The argument to the `-bench` flag selects which benchmarks to run. It's a regular expression matching the names of `Benchmark` functions, with a default value that matches none of them. `.` pattern causes it to match all benchmarks in the `word` package
- The benchmark name's numeric suffix (`8`) indicates the values of `GOMAXPROCS`, which is important for concurrent benchmarks
- Since the benchmark runner initially has no idea how long the operation takes, it makes some initial measurements using small values of `N` and then **extrapolates** to a value large enough for a stable timing measurement to be made
- The reason the loop is implemented by the benchmark function, and not by the calling code in the test driver, is so that the benchmark function has the opportunity to execute any necessary one-time setup code outside the loop without this adding to the measured time of each iteration
    - If the setup code is still perturbing the results, the `testing.B` parameter provides methods to stop, resume, and reset the timer, but these are rarely needed
- An obvious optimization doesn't always yield the expected benefit. The fastest program is often the one that makes the fewest memory allocations (`ch11/word2`)
    - The `-benchmem` will include memory allocation statistics in its report
- Relative timing of 2 different operations
    - If a function takes 1ms to process 1,000 elements, how long will it take to process 10,000 or a million? Such comparisons reveal the asymptotic growth of the running time of the function
    - What's the best size for an I/O buffer? Benchmarks of application throughput over a range of sizes can help us choose the smallest buffer that delivers satisfactory performance
    - Which algorithm performs best for a given job? Benchmarks that evaluate 2 different algorithms on the same input data can often show the strengths and weaknesses of each one on important or representative workloads
- Comparative benchmarks are just regular code. They typically take the form of a single parameterized function, called from several `Benchmark` functions with different values
	
    ```go
    func benchmark(b *testing.B, size int) { /*  */ }
    func Benchmark10(b *testing.B) { benchmark(b, 10) }
    func Benchmark10000(b *testing.B) { benchmark(b, 10000) }
    ```

    - The parameter `size`, which specifies the size of the input, varies across benchmarks but is constant within each benchmark. Resist the temptation to use the parameter `b.N` as the input size. Unless you interpret it as an iteration count for a fixed-size input, the results of your benchmark will be meaningless
- Patterns revealed by comparative benchmarks are particularly useful during program design, but we don't throw the benchmarks away when the program is working. As the program evolves, or its inputs grows, or it's deployed on new operating systems or processors with different characteristics, we can reuse those benchmarks to revisit design decisions
## Profiling
- Benchmarks are useful for measuring the performance of specific operations, but when we're trying to make a slow program faster, we often have no idea where to begin
- Donald Knuth's "Structured Programming with go to Statements"
    > - There's no doubt that the grail of efficiency leads to abuse. Programmers waste enormous amounts of time thinking about, or worrying about, the speed of noncritical parts of their programs, and these attempts at efficiency actually have a strong negative impact when debugging and maintenance are considered. We should forget about small efficiencies, say about 97% of the time: premature optimization is the root of all evil.
    > - Yet we should not pass up our opportunity in that critical 3%. A good programmer will not be lulled into complacency by such reasoning, he will be wise to look carefully at the critical code but only after that code ahs been identified. It is often a mistake to make a priori judgments about what parts of a program are really critical, since the universal experience of programmers who have been using measurement tools has been that their intuitive guesses fail.
- When we wish to look carefully at the speed of our programs, the best technique for identifying the critical code is *profiling*. Profiling is an automated approach to performance measurement based on sampling (采样) a number of profile **events** during execution, then extrapolating from them during a post-process step; the resulting statistical summary is called a profile
- Go supports many kinds of profiling, each concerned with a different aspect of performance, but all of them involve recording a sequence of events of interest, each of which has an accompanying stack trace–the stack of function calls active at the moment of the event
- The `go test` tool has built-in support for several kinds of profiling
- A *CPU profile* identifies the functions whose execution requires the most CPU time
    - The currently running thread on each CPU is interrupted periodically by the operating system every milliseconds, with each interruption recording one profile event before normal execution resumes
- A *heap profile* identifies the statements responsible for allocating the most memory
    - The profiling library samples calls to the internal memory allocation routines so that on average, one profile event is recorded per 512KB of allocated memory
- A *blocking profile* identifies the operations responsible for blocking goroutines the longest, such as system calls, channel sends and receives, and acquisitions of locks
    - The profiling library records an event every time a goroutine is blocked by one of these operations
- Gathering a profile for code under test is as easy as enabling one of the flags below
	
    ```bash
    go test -cpuprofile=cpu.out
    go test -blockprofile=block.out
    go test -memprofile=mem.out
    ```

    - Be careful when using more than 1 flag at a time: the machinery for gathering 1 kind of profile may **skew** the results of others
- It's easy add profiling support to non-test programs too, though the details of how we do that vary between short-lived command-line tools and long-running server applications
    - Profiling is especially useful in long-running applications, so the Go's runtime's profiling features can be enabled under programmer control using the `runtime` API
- Once we've gathered a profile, we need to analyze it using the `pprof` tool. This is part of the Go distribution, but since it's not an everyday tool, it's accessed indirectly using `go tool pprof`
    - It has dozens of features and options, but basic use requires only 2 arguments, the executable that produced the profile and the profile log
    - To make profiling efficient and to save space, the log does not include function names; instead, functions are identified by their addresses. This means that `pprof` needs the executable in order to make sense of the log
        - Although `go test` usually discards the test executable once the test is complete, when profiling is enabled, it saves the executable as `foo.test`, where `foo` is the name of the tested package
- Gather and display s simple CPU profile
	
    ```bash
    go test -run=NONE -bench=ClientServerParallelTLS64 -cpuprofile=cpu.log net/http

    go tool pprof -text -nodecount=10 ./http.test cpu.log
    # shows that elliptic-curve cryptography is important to the performance of this HTTPS benchmark
    ```

    - It's usually better to profile specific benchmarks that have been constructed to be representative of workloads one cares about
    - Benchmarking test cases is almost never representative, which is why we disabled them by using the filter `-run=NONE`
    - The `-text` flag specifies the output format, in this case, a textual table with one row per function, sorted so the hottest functions appear first
    - The `-nodecount=10` limits the results to 10 rows
    - For more subtle problems, you may better off using one of `pprof`'s graphical displays
        - These require [GraphViz](http://www.graphviz.org/)
        - The `-web` flag then renders a directed graph of the functions of the program, annotated by their CPU profile numbers and colored to indicate the hottest functions
- https://blog.golang.org/profiling-go-programs
## Example functions
- The third kind of function treated specially be `go test` is an example function, one whose name starts with `Example`. It has neither parameters nor results
	
    ```go
    func ExampleIsPalindrome() {
        fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
        fmt.Println(IsPalindrome("palindrome"))
        // Output:
        // true
        // false
    }
    ```

- Example functions serve 3 purposes
    1. Documentation (primary)
        - A good example can be a more succinct or intuitive way to convey the behavior of a library function than its prose description, especially when used as a reminder or quick reference
        - An example can also demonstrate the interaction between several types and functions belonging to one API, whereas prose documentation must always be attached to one place, like a type or function declaration or the package as a whole
        - And unlike examples within comments, example functions are real Go code, subject to compile-time checking, so they don't become stale as the code evolves
        - Base on the suffix of the `Example` function, the web-based documentation server `godoc` associates example function with the function or package they exemplify
            - So `ExampleIsPalindrome` would be shown with the documentation for the `IsPalindrome` function, and an example function called just `Example` would be associated with the `word` package as a whole
    2. Examples are executable tests run by `go test`
        - If the example function contains a final `// Output:` comment, the test driver will execute the function and check that what it printed to its standard output matches the text within the comment
    3. An example is hands-on experimentation
        - The `godoc` server at `golang.org` uses the Go Playground to let the user edit and run each example function from within a web browser
        - This is often the fastest way to get a feel for a particular function or language feature