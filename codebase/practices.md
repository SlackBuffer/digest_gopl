
```go
// good
if err := r.ParseForm(); err != nil {
    log.Print(err)
}
// verbose.
err := r.ParseForm()
if err != nil {
    log.Print(err)
}



// normal practice (so that the successful execution path is not indented)
f, err := os.Open(fname)
if err != nil {
    return err
}
f.Stat()
f.Close()
// don't
if f, err := os.Open(name); err != nil {
    return err
} else {
    f.Stat()
    f.Close()
}



// exponential back-off
func WaitForServer(url string) error {
    const timeout= 1 * time.Minute
    deadline := time.Now().Add(timeout)
    for tries := 0; time.Now().Before(deadline); tries++ {
        _, err := http.Head(url)
        if err == nil {
            return nil // success
        }
        log.Printf("server not responding (%s); retrying...", err)
        time.Sleep(time.Second << uint(tries))
    }
    return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}
```