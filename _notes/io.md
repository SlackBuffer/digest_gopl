- Send EOF
    - Unix: `ctrl+d`
    - Windows: `ctrl+z`
- `strings.NewReader`
- `bufio.Scanner`
    - A type that reads input and breaks it into lines or words
    - Often the easiest way to process input that comes naturally in lines

        ```go
        input := bufio.NewScanner(os.Stdin) // or file input
        for input.Scan() {
            fmt.Println(input.Text())
        }
        ```

        - Each call to `input.Scan()` reads the next line and removes the newline character from the end
        - The `Scan` function returns `true` if there's a line and `false` when there's no more input
        - The result can be retrieved by calling `input.Text()`
- `ioutil.ReadFile`
    - Returns a byte slice
    - The byte slice must be converted into a `string` if it needs to be split by `strings.Split`
- > Under the covers, `bufio.Scanner`, `ioutil.ReadFile`, and `ioutil.WriteFile` use the `Read` and `Write` methods of `*os.File`
- `fmt.Fprintf(w io.Writer, a ...interface{}) (n int, err error)`
    - Formats according to a format specifier and writes to `w`
    - Returns the number of bytes written and any write error encountered

        ```go
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v\n", err)
        }
        ```

- `func Sprintf(format string, a ...interface{}) string`
    - Sprintf formats according to a format specifier
    - Returns the resulting string
    - Can be used to convert number into string
- `log.Fatal`
	- `log.Fatal(http.ListenAndServe("localhost:8000", nil))`
- `log.Print`
    - `Print` calls Output to print to the standard logger
    - Arguments are handled in the manner of `fmt.Print`
- Output streams
    - `os.Stdout`, `ioutil.Discard`, `http.ResponseWriter` (all 3 satisfy a common interface `io.Writer`, allowing any of them to be used wherever an output stream is needed)