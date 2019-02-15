# JSON
- `encoding/json`, `encoding/xml`, `encoding/asn1`
- The basic JSON types are numbers (in decimal or scientific notation), booleans, and strings, which are sequences of Unicode code points enclosed in double quotes, with backslash escapes using a similar notation to Go (JSON's `\uhhhh` numeric escapes denote UTF-16 codes, not runes)
- These basic types may be combined recursively using JSON arrays and objects
    - JSON array are used to encode Go arrays and slices
    - JSON objects are used to encode Go maps and structs
- Converting a Go data structure to JSON is called marshaling (`json.Marshal`)

    ```go
    data, err := json.Marshal(movies)
    data, err := json.MarshalIndent(movies, "", " ") // prefix for each line of output; a string for each level of indentation
    ```

    - `Marshal` produces a byte slice containing a very long string with no extraneous white space (hard to read)
    - `json.MarshalIndent` produces neatly indented output
- Marshaling uses the Go struct fields names as the field names for the JSON objects (through reflection)
- **Only exported fields are marshaled**
- A field tag is a string of metadata associated at compile time with the field of a struct
    - A field tag may be any literal string, but it's conventionally interpreted as a space-separated list of `key:"value"` pairs; since they contain double quotation marks, fields tag are usually written with raw string literals
    - The `json` key controls the behavior of the `encoding/json` package, and other `encoding/...` packages follow this convention
    - The first part of the `json` field tag specifies an alternative JSON name for the Go field
    - The additional option `omitempty` indicates that no JSON output should be produces if the fields has the zero value for its type (`false` here) or is otherwise empty
    - Field tags are often used to specify an idiomatic JSON name like `total_count` for a Go field name `TotalCount`

    ```go
    type Movie struct {
        Title string
        Year  int  `json:"released"`
        Color bool `json:"color,omitempty"`
    }
    ```

- The operation of decoding JSON and populating a Go data structure is unmarshaling (`json.Unmarshal`)

    ```go
    var titles []struct{ Title string }
    if err := json.Unmarshal(data, &titles); err != nil {
        log.Fatal("JSON unmarshaling failed: %s", err)
    }
    fmt.Println(titles)
    ```

    - By defining suitable Go data structures, we can select which parts of the JSON input to decode and which to discard
    - When `Unmarshal` returns, it has filled in the slice with the `Title` information; other names in the JSON are ignored
    - The matching process that associates JSON names with Go struct names during unmarshaling is **case-insensitive**, so it's only necessary to use a field tag when there's an underscore in the JSON name but not in the Go name
- The streaming decoder `json.Decoder` allows several JSON entities to be decoded in sequence from the same stream