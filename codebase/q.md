```go
type server struct{}

pb.RegisterGreeterServer(s, &server{})
```