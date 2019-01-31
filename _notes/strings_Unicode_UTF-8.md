# Strings
- A string is an immutable sequence of **bytes**
- Strings may contain arbitrary data, including bytes with value 0
- Text strings are conventionally interpreted as **UTF-8-encoded** sequences of Unicode code points (runes)
- `len` returns the number of bytes (not runes)
- The index operation `s[i]` retrieves the `i`-th byte of string `s`
- The `i`-th byte of a string is not necessarily the `i`-th character of a string, because UTF-8 encoding of a non-ASCII code point requires 2 or more bytes
- The substring operation `s[i:j]` is half-open
    - Either or both of `i` and `j` operands may be omitted, in which case the default values of `0` and `len(s)` are assumed, respectively
- The `+` operator makes a new string by concatenating 2 strings
- Strings may be compared with comparison operators like `==` and `<`
    - The comparison is done **byte by byte**, so the result is the natural lexicographic order
- String **values** are immutable
- Immutability means it's safe for 2 copies of a string to share the same underlying memory, making it cheap to copy strings of any length
- Similarly, a string `s` and a substring like `s[7:]` may safely share the same data, so the substring is also cheap
## String literals
- A string literal is a sequence of **bytes** enclosed in double quotes
- Go source files are always encoded in UTF-8 and Go text strings are conventionally interpreted as UTF-8, we can **include Unicode points in string literals**
- Within a double-quoted string literal, escape sequences that begin with a backslash `\` can be used to insert arbitrary byte values into the string
- Arbitrary bytes can also be included in literal strings using hexadecimal or octal escapes
    - A hexadecimal escape is written `\xhh`, with exactly 2 hexadecimal digits `h` (or uppercase)
    - An octal escape is written `\ooo` with exactly 3 octal digits `o` (0 through 7) not exceeding `\377` (255 in decimal)
    - Both denotes a single byte with the specified value
- A **raw string** uses `...`
- Within a raw string literal, no escape sequences are processed; the contents are taken literally, including backslashes or newlines
    - So a raw string literal may spread over several lines in the program source
    - The only processing is that **carriage returns are deleted** so that the value of the string is the same on all platforms, including those that conventionally put carriage returns in text files
- Raw string literals are useful for regular expressions, HTML templates, JSON literals, command usage messages, and the like, which often extend over multiple lines
## Unicode
- ASCII, or more precisely US-ASCII, uses 7 bits to represent 128 characters
- Unicode collects all of the characters in all of the world's writing systems, plus accents and other diacritical mark, control codes like tab and carriage return, and plenty of esoterica, and **assigns each one a standard number called Unicode code point**, or in Go's terminology, a rune
- The natural data type to hold a single rune is `int32`, and that's what Go uses; it has the synonym `rune` for precisely this purpose
- UTF-32 or UCS-4 represent a sequence of runes as s sequence of `int32` values. The encoding of each Unicode code point has the same size, 32 bits
    - This is simple and uniform
    - Most computer-readable text is in ASCII, which requires only 8 bits or 1 byte per character
    - All the characters in widespread use still number fewer that 65536, which would fit in 16 bits
## UTF-8
- UTF-8 is a variable-length encoding of Unicode points as bytes
- It uses between 1 and 4 bytes to represent each rune, but only 1 byte for ASCII character, and only 2 or 3 bytes for most runes in common use
- The high-order bits of the first byte of the encoding for a rune indicate how many bytes follow
   1. A high-order `0` indicates 7-bit ASCII where each rune takes 1 byte (identical to conventional ASCII)
   2. A high-order `110` indicates that the rune takes 2 bytes; the second byte begins with `10`

        ```
        # x 的个数，即决定数值范围的 bit 数
        # 2^7=128, 2^11=2048, 2^16=65536, 2^21=2097152 (larger than 1114111(0x10ffff))

        0xxxxxxx                             runes 0−127     (ASCII) 1 byte
        110xxxxx 10xxxxxx                    128−2047        (values less than 128 unused) 2 bytes
        1110xxxx 10xxxxxx 10xxxxxx           2048−65535      (values less than 2048 unused) 3 bytes
        11110xxx 10xxxxxx 10xxxxxx 10xxxxxx  65536−0x10ffff  (other values unused) 4 bytes
        ```