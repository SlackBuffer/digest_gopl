- [ ] graceful shutdown
- [ ] http handler recover from panic
- [ ] `go test` 一次通过测试后，再运行测试命令还是得到前一次测试通过的结果，需要重新编辑并保存测试代码才会重新测试
- [x] go run: cannot run non-main package
- [x] curl and browser behave differently (Ex1.12.1): additional request '/favicon.ico'
- [x] Unnamed type

    ```go
    // unnamed type `struct{ I int }`
    var x struct{ I int }

    // named type `Foo`
    type Foo struct{ I int }
    var y Foo
    ```

    - Composite types are known as unnamed types, as they use a **type literal** to represent the structural definition of the type, instead of using a simple name identifier, Unlike its named counterpart, unnamed composite types use literals for value initialization that are composed of type (itself) and a literal text that represents the value
    - > https://medium.com/learning-the-go-programming-language/types-in-the-go-programming-language-65e945d0a692
    - > https://stackoverflow.com/questions/32983546/named-and-unnamed-types
- [ ] So a declaration may refer to itself (like functions, types!) or to another that follows it, letting us declare recursive or [ ] mutually recursive types and functions
- [ ] `^`: bitwise XOR (bitwise exclusive OR); `&^`: bit clear (AND NOT)
- [x] Print �: `fmt.Println(string(1234567))`
- [ ] object files (`.a`?)
- buffer
    - https://stackoverflow.com/questions/47186741/how-to-get-the-json-from-the-body-of-a-request-on-go/47295689#47295689
    - https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
    - https://github.com/gin-gonic/gin/issues/1295#issuecomment-402254995