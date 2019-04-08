# A tour of Go
- A goroutine is a lightweight thread managed by the Go runtime
- `go f(x, y, z)` starts a new goroutine running `f(x, y, z)`
- The evaluation of `f`, `x`, `y`, `z`` happens in the current goroutine and the execution of `f` happens in the new goroutine
- Goroutines run in the **same address space**, so access to shared memory must be synchronized
- Channels are a **typed conduit** through which you can send and receive values with the channel operator `<-`

    ```go
    ch <- v // send v to channel ch
    v := <-ch // receive from ch, and assign value to v
    ```

    - The data flows in the direction of the arrow
- Like maps and slices, channels must be created before use

    ```go
    ch := make(chan int)
    ```

- By default, sends and receives block until the other side is ready. This allows goroutines to synchronize without explicit locks or condition variables

    ```go
    func sum(s []int, c chan int) {
        sum := 0
        for _, v := range s {
            sum += v
        }
        c <- sum
    }
    func main() {
        s := []int{7, 2, 8, -9, 4, 0}
        c := make(chan int)
        go sum(s[:len(s) / 2], c)
        go sum(s[len(s) / 2:], c)
        x, y := <-c, <-c
        fmt.Println(x, y, x+y)
    }
    ```

- Channels can be buffered

    ```go
    // provide the buffer length as the second parameter
    ch := make(chan int, 100)
    ```

- Sends to a buffered channel block only when the buffer is full
- Receivers block when the buffer is empty
- A sender can `close` a channel to indicate that no more values will be sent. Receiver can test whether a channel has been closed by assigning a second parameter to the receiver expression

    ```go
    v, ok := <-ch
    ```

    - `ok` is `false` if there are no more values to receive and the channel is closed
- Only the sender should close a channel, never the receiver
- Sending on a closed channel will cause a panic
- The loop `for i := range c` receive values from the channel repeatedly until it is **closed**
- Channels aren't like files; you don't usually need to close them
- Closing is only necessary when the receiver must be told there're no more values coming, such as to terminate a `range` loop

    ```go
    func fibonacci(n int, c chan int) {
        x, y := 0, 1
        for i := 0; i < n; i++ {
            c <- x
            x, y = y, x+y
        }
        close(c)
    }
    func main() {
        c := make(chan int, 10)
        go fibonacci(cap(c), c)
        for i := range c {
            fmt.Println(i)
        }
    }
    ```

- The `select` statement lets a goroutine wait on multiple communication operations
- A `select` blocks until one of its cases can run, then it executes that case
- It chooses one at random if multiple are ready

    ```go
    func fibonacci(c, quit chan int) {
        x, y := 0, 1
        for {
            // 协程间交互
            select {
            // 1. 写入 channel c，阻塞，直至 c 被读取
            case c <- x:
                x, y = y, x+y
            case <-quit:
                fmt.Println("quit")
                return
            }
        }
    }
    func main() {
        c := make(chan int)
        quit := make(chan int)
        go func() {
            for i := 0; i < 10; i++ {
                // 从 2. channel 中读取；另一协程中的无限循环写入操作得以继续，直至 return
                fmt.Println(<-c)
            }
            quit <- 0
        }()
        fibonacci(c, quit)
    }
    ```
- The `default` case in a `select` is run if no other case is ready
- Use a `default` case to try and send or receive without blocking

    ```go
    func main() {
        tick := time.Tick(100 * time.Millisecond)
        boom := time.After(500 * time.Millisecond)
        for {
            select {
            case <-tick:
                fmt.Println("tick.")
            case <-boom:
                fmt.Println("BOOM!")
                return
            default:
                fmt.Println("   .")
                time.Sleep(50 * time.Millisecond)
            }
        }
    }
    ```

- [ ] https://tour.golang.org/concurrency/8?a=1
- The concept of only one goroutine can access a variable at a time to avoid conflicts is called mutual exclusion. The conventional name for the data structure that provides it is mutex
- Go's standard library provides mutual exclusion with `sync.Mutex` and its 2 methods `Lock` and `Unlock`

    ```go
    type SafeCounter struct {
        v map[string]int
        mux sync.Mutex
    }
    func (c *SafeCounter) Inc(key string) {
        c.mux.Lock()
        // lock so only one goroutine at a time can access the map c.v
        c.v[key]++
        c.mux.Unlock()
    }
    func (c *SafeCounter) Value(key string) int {
        // lock so that only one goroutine at a time can access the map c.v
        c.mux.Lock()
        // use defer to make sure the mutex wil be unlocked
        defer c.mux.Unlock()
        return c.v[key]
    }

    c := SafeCounter{v : make(map[string]int)}
    for i := 0; i < 1000; i++ {
        go c.Inc("somekey")
    }
    time.Sleep(time.Second)
    fmt.Println(c.Value("someky"))
    ```

# TLDR
- A channel is a **communication mechanism** that allows one goroutine to pass values of a specified type to another goroutine
- `main` runs in a goroutine and the `go` statement creates additional goroutines
- When one goroutine attempts a send or receive **on a channel**, it blocks until another goroutine attempts the corresponding receive or send operation, at which point the value is transferred and both goroutines proceed
- Communicating sequential process (CSP) - a model of concurrency in which values are passed between independent activities (goroutines) but variables are for most part confined to a single activity
# Goroutines
- In Go, each concurrently executing activity is called a goroutine
- When a program starts, its only goroutine is the one that calls the `main` function (called **main goroutine**)
- New goroutine are created by the `go` statement (new goroutine calls other functions)
    - Syntactically, a `go` statement is an ordinary function or method call prefixed by the keyword `go`
    - A `go` statement causes the function to be called in a newly created goroutine
- The `go` keyword itself completes immediately

    ```go
    f() // wait for f() to return
    go f() // create a new goroutine that calls f(); won't wait
    ```

- When the `main` function returns, **all** goroutines are abruptly terminated and the **program exists**
    - Other than by returning from `main` or existing the program, there is no programmatic way for one goroutine to stop another
    - There are ways to communicate with a goroutine to request that it stop itself
## Example: concurrent clock server
- `time.Time.Format` provides a way to format date and time information by example
    - Its argument is a template indicating how to format a reference time, specially `Mon Jan 2 03:04:05PM 2006 UTC-0700`
    - The reference time has 8 components. Any collection of them can appear in the `Format` string in **any order** and in **a number of formats**
    - The selected components of will be displayed in the selected formats
    - The `time` package defines templates for many standard time formats， such as `time.RFC1123`
    - The same mechanism is used in reverse when passing a time using `time.Parse`
## Example: concurrent echo server
- In adding `go` keywords, we had to consider carefully that it's safe to call methods of `net.Conn` concurrently, which is not true for most types
# Channels
- A channel is a communication mechanism that lets one goroutine send values to another goroutine
- Each channel is a conduit (pipe) for values of a particular type, called the channel's **element type**
    - The type of a channel whose elements have type `int` is written `chan int`
- Use `make` to create a channel

    ```go
    ch := make(chan int)
    ```