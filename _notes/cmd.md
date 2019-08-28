<!-- - Add a `&` to the end of of a command to run it in the background -->
- `lsof -i:8000`
<!-- - `kill -9 pid` -->
- `killall clock1` kills all processes with the given name
- Alternative to `go get`

    ```bash
    mkdir -p $GOPATH/src/golang.org/x/
    git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
    go install net
    ```

- 当前终端走代理：`export https_proxy=socks5://127.0.0.1:1086`
    - ShadowSocks - Local socks5 listen address && port