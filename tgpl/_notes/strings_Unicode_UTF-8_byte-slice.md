# Strings
- A string is an **immutable sequence of bytes**
- Strings may contain arbitrary data, including bytes with value 0
- Text strings are conventionally interpreted as **UTF-8-encoded** sequences of Unicode code points (runes)
- `len` returns the number of bytes (**not runes**)
- The index operation `s[i]` retrieves the `i`-th **byte** of string `s`
    - The `i`-th byte of a string is not necessarily the `i`-th character of a string, because UTF-8 encoding of a non-ASCII code point requires 2 or more bytes
- The substring operation `s[i:j]` is half-open
    - Either or both of `i` and `j` operands may be omitted, in which case the default values of `0` and `len(s)` are assumed, respectively
- The `+` operator makes a new string by concatenating 2 strings
- Strings may be compared with comparison operators like `==` and `<`
    - The comparison is done **byte by byte**, so the result is the natural lexicographic order
- String **values** are immutable
    - Immutability means it's safe for 2 copies of a string to **share the same underlying memory**, making it cheap to copy strings of any length
    - Similarly, a string `s` and a substring like `s[7:]` may safely share the same data, so the substring is also cheap
## String literals
- A string literal is a sequence of bytes enclosed in double quotes
- Go source files are always encoded in UTF-8 and Go text strings are conventionally interpreted as UTF-8, we can **include Unicode points in string literals**
- Within a double-quoted string literal, escape sequences that begin with a backslash `\` can be used to insert arbitrary byte values into the string
    - One set of escapes handles ASCII control codes like newline, carriage return, and tab
- Arbitrary bytes can also be included in literal strings using hexadecimal or octal escapes
    - A hexadecimal escape is written `\xhh`, with exactly 2 hexadecimal digits `h` (or uppercase)
    - An octal escape is written `\ooo` with exactly 3 octal digits `o` (0 through 7) not exceeding `\377` (255 in decimal)
    - Both denotes a single byte with the specified value
- A **raw string** is written `...`, using backquotes instead of double quotes
- Within a raw string literal, no escape sequences are processed; the contents are taken literally, including backslashes or newlines
    - So a raw string literal may spread over several lines in the program source
    - The only processing is that **carriage returns are deleted** so that the value of the string is the same on all platforms, including those that conventionally put carriage returns in text files
- Raw string literals are useful for regular expressions, HTML templates, JSON literals, command usage messages, and the like, which often extend over multiple lines
## Unicode
- ASCII, or more precisely US-ASCII, uses 7 bits to represent 128 characters
- Unicode collects all of the characters in all of the world's writing systems, plus accents and other diacritical mark, control codes like tab and carriage return, and plenty of esoterica, and **assigns each one a standard number called Unicode code point**, or in Go's terminology, a *rune*
- The natural data type to hold a single rune is `int32`, and that's what Go uses; it has the synonym `rune` for precisely this purpose
- UTF-32 or UCS-4 represent a sequence of runes as a sequence of `int32` values. The encoding of each Unicode code point has the same size, 32 bits
    - This is simple and uniform
    - Most computer-readable text is in ASCII, which requires only 8 bits or 1 byte per character
    - All the characters in widespread use still number fewer that 65536, which would fit in 16 bits
## UTF-8
- UTF-8 is a *variable-length* encoding of Unicode points as bytes. It uses between 1 and 4 bytes to represent each rune, but only 1 byte for ASCII character, and only 2 or 3 bytes for most runes in common use
- The high-order bits of the first **byte** of the encoding for a rune indicate how many bytes follow (高位 1 的个数决定编码占用的字节数)
   1. A high-order `0` indicates 7-bit ASCII where each rune takes 1 byte (identical to conventional ASCII)
   2. A high-order `110` indicates that the rune takes 2 bytes; the second byte begins with `10`

        ```
        # x 的个数，是决定数值范围的 bit 数
        # 2^7=128, 2^11=2048, 2^16=65536, 2^21=2097152 (larger than 1114111(0x10ffff))

        0xxxxxxx                             runes 0−127     (ASCII) 1 byte
        110xxxxx 10xxxxxx                    128−2047        (values less than 128 unused) 2 bytes
        1110xxxx 10xxxxxx 10xxxxxx           2048−65535      (values less than 2048 unused) 3 bytes
        11110xxx 10xxxxxx 10xxxxxx 10xxxxxx  65536−0x10ffff  (other values unused) 4 bytes
        ```

- A variable-length encoding precludes direct indexing to access the `n`-th character of a string
- Advantages
    - The encoding is compact, compatible with ASCII, and self-synchronizing: its' possible to find the beginning of a character by backing up no more than 3 bytes
    - It's a prefix code, so it can be decoded from **left to right** without any ambiguity or lookahead
    - No rune's encoding is a substring of any other, or even a sequence of others, so you can search for a rune just by searching for its bytes, without worrying about the preceding context
    - The lexicographic byte order equals the Unicode code point order, so sorting UTF-8 works naturally
    - There are no embedded NUL (zero) bytes, which is convenient for programming languages that use NUL to terminate strings
- Packages: `unicode`, `unicode/utf8`
- Many Unicode characters are hard to type on a keyboard or to distinguish visually from similar-looking ones; some are invisible
- **Unicode escapes** in Go string literals allow us to specify them by their **numeric code point value**. 2 forms
   1. `\uhhhh` for a 16-bit value
   2. `\Uhhhhhhhh` for a 32-bit value

        ```go
        // 4 string literals of valid UTF-8 encoding of representing the same thing
        "世界"
        "\xe4\xb8\x96\xe7\x95\x8c"
        // 1110-0100 10-111000 10-010110
        // 0100111000010110 - 4e16
        "\u4e16\u754c"
        "\U00004e16\U0000754c"
        ```

        - Each `h` is a hexadecimal value
- Unicode escapes may also be used in rune literals

    ```bash
    # these 3  are equivalent
    '世' '\u4e16' '\U00004e16'
    ```

    - A rune whose value is less than 256 may be written with a single hexadecimal escape, such as `\x41` for `A`, but for higher values, a `\u` or `\U` must be used. Consequently, `'\xe4\xb8\x96'` is not a legal rune literal (larger than 256), even though those 3 bytes are a valid are a valid UTF-8 encoding for a single code point
- Count runes

    ```go
    fmt.Println(len("Hello, 世界")) // 13
    fmt.Println(utf8.RuneCountInString("Hello, 世界")) // 9

    n := 0
    for _, _ = range "Hello, 世界" {
        n++
    } // 9

    n := 0
    for range "Hello, 世界" {
        n++
    } // 9
    ```

- Process individual Unicode character

    ```go
    s := "Hello, 世界"
    // method 1
    for i := 0; i < len(s); {
        r, size := utf8.DecodeRuneInString(s[i:])
        fmt.Printf("%d\t%c\n", i, r)
        i += size
    }

    // method 2
    for i, r := range s {
        fmt.Printf("%d\t%q\t%d\n", i, r, r)
    }
    ```

    - Go's `range` loop, when applied to a string, **performs UTF-8 decoding implicitly**
    - It's mostly a matter of convention in Go that text strings are interpreted as UTF-8-encoded sequences of Unicode code points, but for correct use of `range` loops on strings, it's more than a convention, it's a necessity
- Each time a UTF-8 decoder, whether explicit in a call to `utf8.DecodeRuneInString` or implicit in a `range` loop, consumes an **unexpected input byte**, it generates a special Unicode replacement character, `'\uFFFD'` (�, a white question mark inside a black hexagonal or diamond-like shape). Like when range over a string containing arbitrary binary data, or for that matter, UTF-8 data containing errors
    - When a program encounters this rune value, it's often a sign that some upstream part of the system that generated the string data has been careless in its treatment of text encodings

    <!-- ```go
    s := "\u4e16\u754c\134"
    for _, r := range s {
        fmt.Printf("%c\t", r)
    }
    ``` -->

- UTF-8 is exceptionally convenient as an interchange format **but** within a program, runes may be more convenient because they are of uniform size and thus easily indexed in arrays and slices
- A `[]rune` conversion applied to a UTF-8-encoded string returns the sequence of Unicode code points that the string encodes

    ```go
    s := "世界"
    fmt.Printf("% x\n", s)  // e4 b8 96 e7 95 8c
    r := []rune(s)
    fmt.Printf("%x\n", r) // [4e16 754c]

    fmt.Println(string(r)) // 世界
    ```

    - Converting a slice of runes to a string produces the **concatenation** of the UTF-8 encoding of each rune
- **Converting an integer value to a string interprets the integer as a rune value**, and yields the UTF-8 representation of that rune

    ```go
    fmt.Println(string(65)) // A
    fmt.Println(string(0x4eac)) // 京

    fmt.Println(string(1234567)) // �
    ```

    - If the rune is invalid, the replacement character is substituted
## Strings and byte slices
- `bytes`, `strings`, `strconv`, `unicode`
- Strings are immutable, **building up strings incrementally** can involve a lot of allocation and copying. In such cases, it's more efficient to use the `bytes.Buffer` type
- `path`, `path/filepath`
- Strings can be converted to byte slices and back again

    ```go
    s := "abc"
    b := []byte(s)
    s2 := string(b)
    ```    

- `[]byte(s)` allocates a new byte array holding a copy of the bytes of `s`, and yields a slice that references the entirety of that array
    - An optimized compiler may be able to avoid the allocation and copying in some cases, but in general copying is required to ensure that the bytes of `s` remain unchanged even if those of `b` are subsequently modified
- The conversion from byte slice to string also makes a copy, to **ensure immutability** of the resulting string `s2`
- To avoid conversions and unnecessary memory allocation, many of the utility functions in the `bytes` package directly parallel their counterparts in the `string` package
- The `bytes` package provides the `Buffer` type for the efficient manipulation of byte slices. A `Buffer` starts out empty but grows as data of types like `string`, `byte`, and `[]byte` are written to it
    - When appending the UTF-8 encoding of an arbitrary rune to a `bytes.Buffer`, it's best to use `bytes.Buffer`'s `WriteRune` method, but `WriteByte` is fine for ASCII characters
## Conversions between strings and numbers
- `strconv`
- Convert an integer to a string

    ```go
    x := 123
    y := fmt.Sprintf("%d", x) // 1
    fmt.Println(y, strconv.Itoa(x)) //2

    fmt.Println(strconv.FormatInt(int64(x), 2))
    ```

    - `FormatInt` and `FormatUint` can be used to format numbers in a different base
- The `fmt.Printf` verbs `%b`, `%d`, `%o`, and `%x` are often more convenient than `Format` functions, especially if we want to include additional information besides the number

    ```go
    s := fmt.Sprintf("x=%b", x)
    ```

- Parse a string representing an integer

    ```go
    x, err := strconv.Atoi("123")
    y, err := strconv.ParseInt("123", 10, 64) // base 10, up to 64 bits
    ```

    - 0 as the third argument implies `int`
- Sometimes `fmt.Scanf` is useful for parsing input that consists of orderly mixtures of strings and number all in a single line, but it can be inflexible