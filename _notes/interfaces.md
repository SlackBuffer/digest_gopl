- Interfaces are **abstract types** that allow us to treat different **concrete types** in the same way based on what **methods** they have, not how they are represented or implemented
- Go's interfaces are **satisfied implicitly**
    - There's no need to declare all the interfaces that a given concrete type satisfies; simply **possessing the necessary methods is enough**
- This design lets you **create new interfaces** that are **satisfied by existing concrete types [ ] without changing the existing types**, which is particularly useful for types defined in packages that you don't control
# Interface as contracts
- A concrete type specifies the exact representation of its values and exposes the intrinsic operations of that representation
    - Such as arithmetic for numbers, or indexing, `append` and `range` for slices
- A concrete type may also provide additional behaviors through its methods
- When you have a value of a concrete, you know exactly what it is and what you can do with it
- An interface is an abstract type
    - It doesn't expose the representation or internal structure of its values, or the set of basic operations they support
    - It reveals only **some of their methods**
- When you have a value of an interface, you know nothing about what it is; you know only what it can do, or more precisely, what behaviors are provided by its methods
- `fmt.Fprintf`

    ```go
    package fmt

    // Fprintf formats according to a format specifier and writes to w
    func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)

    func Printf(format string, args ...interface{}) (int, error) {
        return Fprintf(os.Stdout, format, args...)
    }

    func Sprintf(format string, args ...interface{}) string {
        var buf bytes.Buffer
        Fprintf(&buf, format, args...)
        return buf.String()
    }
    ```

    - The `F` prefix of `Fprintf` stands for file and indicates that the formatted output should be written to the file provided as the first argument
    - In the `Printf` case, `os.Stdout` is an `*os.File` (a file)
    - In the `Sprintf` case, the argument is not a file, though is superficially resembles one: `&buf` is a pointer to a memory buffer to which bytes can be written
    - The first parameter of `Fprintf` is not a file. It's an `io.Writer` (interface type)

        ```go
        package io
        // `Writer` is the interface that wraps the basic Write method
        type Writer interface {
            // Write writes `len(p)` bytes from p to the underlying data stream
            // It returns the number of bytes written from p and any error
            // encountered that caused the write to stop early
            // Must return a non-nil error if it returns n < len(p)
            // Must not modify the slice data, even temporarily
            Write(p []byte) (n int, err error)
        }
        ```

- The `io.Writer` interface defines the **contract** between `Fprintf` and its callers
   1. The contract **requires** that the caller provides a value of a concrete type like `*os.File` or `*bytes.Buffer` that has a method called `Write` with the appropriate signature and behavior
   2. The contract **guarantees** that `Fprintf` will do its job given any value that satisfies the `io.Writer` interface
- `Fprintf` may not assume that it is writing to a file or to memory, only that it **can call** `Write`
    - Because `fmt.Fprintf` assumes nothing about the representation of the value (of its first argument) and relies only on the behaviors guaranteed by the `io.Writer` contract, we can safely pass a value of any concrete type that satisfies `io.Writer` as the first argument to `fmt.Fprintf`
- The freedom to substitute one type for another that satisfies the same interface is called **substitutability**
- `fmt.Stringer`

    ```go
    package fmt
    // The String method is used to print values passed as an operand to any format that accepts a string
    // or to an unformatted printer such as Print
    type Stringer interface {
        String() string
    }
    ```

# Interface type
- An interface type specifies a set of methods that a concrete type must possess to be considered an **instance** of that interface
    - `io.Writer` interface provides an abstraction of all the types to which bytes can be written, which includes files, memory buffers, network connections, HTTP clients, archivers, hashers, and so on
    - A `Reader` interface represents any type from which you can read bytes
    - A `Closer` interface is any value you can close, such as a file or a network connection

    ```go
    package io
    type Reader interface {
        Reader(p []byte) (n int, err error)
    }
    type Closer interface {
        Close() error
    }
    ```

- Declarations of new interface types as combinations of existing ones

    ```go
    type ReaderWriter interface {
        Reader
        Writer
    }
    ```

    - This syntax lets us name another interface as a shorthand for writing out all of its methods. Called embedding an interface
- These 2 declaration style can be mixed

    ```go
    type ReaderWriter interface {
        Reader(p []byte) (n int, err error)
        Writer
    }
    ```

- The order in which the methods appear don't matter
# Interface satisfaction
- **A type satisfies an interface** if it possesses all the methods the interface requires
    - `*os.File` satisfies `io.Reader`, `Writer`, `Closer`, `ReadWriter`
    - `*bytes.Buffer` satisfies `Reader`, `Writer`, `ReadWriter`
- As a shorthand, Go programmers often say that a concrete type **"is a"** particular interface, meaning that it satisfies the interface
    - A `*bytes.Buffer` is an `io.Writer`; an `*os.File` is an `io.ReadWriter`
- An expression may be assigned to an interface if its type satisfies the interface
    - This rule applies even when the right-hand side is itself an interface
- It's legal to call a `*T` method on an argument of type `T` so long as the argument is a variable; the compiler implicitly takes its address. This is mere **syntactic sugar**
- A value of type `T` does not possess all the methods that a `*T` pointer does, and as a result it might satisfy fewer interfaces

    ```go
    type IntSet struct { /* ... */ }
    func (*IntSet) String() string
    // cannot call that method on a non-addressable IntSet value
    var _ = IntSet{}.String() // compile error
    // s is a variable and &s has a String method
    var s IntSet
    var _ = s.String()
    var _ fmt.Stringer = &s
    var _ fmt.Stringer = s // compile error: IntSet lacks String method
    ```

- An interface wraps and conceals the concrete type and value that it holds. Only the methods revealed by the interface type may be called, even if the concrete type has others

    ```go
    os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
    os.Stdout.Close // OK: *os.File has Close method

    var w io.Writer
    w = os.Stdout
    w.Write([]byte("hello")) // OK: *os.File has Write method
    w.Close // compile error: io.Writer lacks Close method
    ```

- The type `interface{}` (empty interface type) places no demands on the types that satisfy (implement) it, we can assign any value to the empty interface
- Since interface satisfaction depends only on the methods of the 2 types involved, there's no need to declare the relationship between a concrete type and the interfaces it satisfies
    - That said, it's occasionally useful to document and **assert** the relationship when it's intended but not otherwise enforced by the program

    ```go
    // asserts at compile time that a value of type *bytes.Buffer satisfies `io.Writer`
    var w io.Writer = new(bytes.Buffer)
    ```

    - We needn't allocate a new variable since any value of type `*bytes.Buffer` will do, even `nil`, which we write as `(*bytes.Buffer)(nil)` using an explicit conversion
    - Since we never intend to refer to `w`, we can replace it with the blank identifier

        ```go
        var _ io.Writer = (*bytes.Buffer)(nil)
        ```

- Non-empty interface types such as `io.Writer` are most often satisfied by a **pointer** type, particularly when one or more of the interface methods implies some kind of **mutation to the receiver**, as the `Writer` method does
    - A pointer to a struct is an especially common method-bearing type
- Pointer types are by not the only types that satisfy interfaces
- Even interfaces with mutator methods may be satisfied by one of Go's other reference types
- Basic types may satisfy interfaces
- A concrete type may satisfy many unrelated interfaces
- Interfaces are one useful way to group related concrete types together and express the facets they share in common

    ```go
    type Audio interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string // e.g., MP3, WAV
    }
    type Video interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string // e.g., MP3, WMV
        Resolution() (x, y int)
    }
    type Streamer interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string // e.g., MP3, WAV
    }
    ```

    - In order to handle `Audio` and `Video` items in the same way, we can define a `Streamer` interface to represent their common aspects without changing any existing type declarations
- Each grouping of concrete types based on their shared behaviors can be expressed as an interface type
- Unlike class-based languages, in which the set of interfaces satisfied by a class is explicit, in Go,we can define new abstractions or groupings of interest when we need them, without modifying the declaration of concrete type
# Parsing flags with `flag.Value`
- `flag.Value` is the interface to the value stored in a flag

    ```go
    package flag
    type Value interface {
        String() string
        Set(string) error
    }
    ```

    - `String` formats the flag's value for use in command-line help messages; every `flag.Value` is also a `fmt.Stringer`
    - `Set` parses its string argument and updates the flag value
- In order to define new flag notations for our own ata types, we need only define a type that satisfies the `flag.Value` interface