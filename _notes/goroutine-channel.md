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

- As with maps, a channel is a **reference** to the data structure created by `make`
- The zero value of a channel is `nil`
- 2 channels of the same type may be compared using `==`
    - The comparison is true if both are references to the same channel data structure
    - A channel may be compared to `nil`
- A channel has 2 principal operations, send and receive, collectively known as communications
- A send statement transmits a value from one goroutine, through the channel, to another goroutine executing **a corresponding receive operation**
    - In a send statement, the `<-` separates the channel and value operands
    - In a receive expression, `<-` precedes the channel operand
    - A receive expression whose result is not used is a valid statement

    ```go
    ch <- x // send
    x = <-ch // receive
    <-ch // receive; result is discarded
    ```

- Channels support a third operation, `close`, which sets a flag indicating that no more values will ever be send on this channel; subsequent attempts to send will panic

    ```go
    // built-in `close`
    close(ch)
    ```

- Receive operations on a closed channel yield the values that have been send until no more values are left; any receive operations **thereafter** complete immediately and yield the **zero value of the channel's element type**
- A channel created with a simple call to `make` is called an unbuffered channel
- `make` accepts an optional second argument, an integer called the channel's **capacity**. If the capacity is non-zero, `make` creates a buffered channel
## Unbuffered channels
- A send operation on an unbuffered channel **blocks the sending goroutine** until **another goroutine** executes a corresponding receive **on the same channel**, at which point the value is transmitted and both goroutines may continue
- Conversely, if the receive operation was attempted first, the receiving goroutine is block until another goroutine performs a send on the same channel
- Communication over an unbuffered channel (synchronous channel) causes the sending and receiving goroutines to **synchronize**
- When a value is sent on an unbuffered channel, the receipt of the value **happens before** the reawakening of **sending** goroutine
- In a discussion of concurrency, when we say `x` happens before `y`, we don't mean merely that `x` occurs earlier in time than `y`; we mean that it's **guaranteed** to do so and that all its prior effects, such as updates to variables, are complete and that you **may rely on them**
- When `x` neither happens before `y` nor after `y`, we say `x` is concurrent with `y`
- Messages sent over channels have 2 important aspects
   1. Each message has a value
   2. Sometimes the fact of communication and the moment at which it occurs are just important
        - Messages are called events when this aspect is to be stressed
        - When the event carries no additional information and its sole purpose is synchronization, we'll emphasize this by using a channel whose element type is `struct{}`
            - It's common to use a channel of `bool` or `int` for the same purpose since `done <- 1` is shorter that `done <- struct{}{}`
## Pipelines
- Channels can be used to connect goroutines together so that the output of one is the input to another. This is called a pipeline
- If the sender knows that no further values will ever be sent on a channel, it's useful to communicate this fact to the receiver goroutine so that they stop waiting
- This is accomplished by closing the channel using the build-in `close` function
- After a channel has been closed, any further send operations on it will **panic**
- After the last sent element has been received, all subsequent receive operations will proceed **without blocking** but will yield a zero value of channel's element type 
- A variant of the receive operation produces 2 results: the received channel element, plus a boolean value, which is `true` for a successful receive and `false` for a receive on a closed and drained channel
- Using a `range` loop to iterate over channels is a more convenient syntax for receiving all values sent on a channel and terminating the loop after the last one
- Needn't close every channel when you've finished with it
- It's only necessary to close a channel when it's important to tell the receiving goroutines that all data have been sent
- A channel that the garbage collector determines to be unreachable will have its resources reclaimed whether or not it is closed
    - Don't confuse this with the close operation for open files. It's important to call the `Close` method on every file when you've finished with it
- Attempting to close an **already-closed** channel causes a panic, as does closing a **nil** channel
- Closing channels has another use as a **broadcast** mechanism
## Unidirectional channel type