
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
```