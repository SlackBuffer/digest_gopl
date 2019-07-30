```go
// [x]
type server struct{}
// see pointer snippets
pb.RegisterGreeterServer(s, &server{})
```