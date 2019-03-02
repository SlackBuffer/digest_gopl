- Add a `&` to the end of of a command to run it in the background
- `lsof -i:8000`, `kill -9 pid`
- Alternative to `go get`

    ```bash
    mkdir -p $GOPATH/src/golang.org/x/
    git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
    go install net
    ```