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