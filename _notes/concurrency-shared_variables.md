# Sum
- In the concurrency-safe one, `icons` variable is **assigned during package initialization**, which *happens before* the program's `main` function starts running. Once initialized, `icons` is never modified
    - **Initialization proceeds from the bottom up; the `main` package is the last to be initialized**
- By convention, the variables guarded by a mutex are declared immediately after the declaration of the mutex itself
- 加锁意味着 critical section 里的操作在某个时刻只能由一个协程进行
    - 可以用 unbuffered channel 来理解
- Non-blocking cache
    - 构造一个数据结构，只在第一次请求时判空成立，该请求负责去做实际的请求，并在拿回结果请写入前阻塞其它读取操作，写入成功后广播消息；其它协程的请求过来，发现判空不成立，无需再去请求，等到收到写入成功的广播消息后去完成读取
        - `ch9/memo4`: Allow may goroutines to access the variable, but only one at a time (shared variables and locks, *mutual exclusion*)
        - `ch9/memo5`: Variable *confined* to a single goroutine (communicating sequential process)
- The Go mantra: **"Do not communicate by sharing memory; instead, share memory by communicating"**
# Race conditions
- In a program with 2 or more goroutines, the steps within each goroutine happen in the familiar order, but in general we don't know whether an event `x` in one goroutine happens before an event `y` in another goroutine, or happens after it, or is simultaneous with it. When we cannot confidently say that one event happens before the other, then the events `x` and `y` are *concurrent*
- Consider a function that works correctly in a sequential program. A function is *concurrency-safe* if it continues to work correctly even when called concurrently, that is, from 2 or more goroutines with no additional synchronization
- A type is concurrency-safe if all its accessible methods and operations are concurrency-safe
- We can make a program concurrency-safe without making every concrete type in that program concurrency-safe. Concurrency-safe types are the **exceptions** rather than the rule. You should access a variable concurrently only if the documentation for its type says that it's safe
- We avoid concurrent access to most variables either by *confining* them to a single goroutine or by maintaining a higher-level invariant of *mutual exclusion*
- In contrast, exported package-level functions are generally expected to be concurrency-safe. Since package-level variables cannot be confined to a single goroutine, functions that modify them must enforce mutual exclusion
- There are many reasons a function might not work when called concurrently, including deadlock, livelock, and resource starvation
    - > [deadlock, starvation](http://quicktechie.com/cs/all-articles/65-java-deadlock-starvation-and-livelock)
    - > [livelock](https://www.techopedia.com/definition/3723/livelock)
- A race condition is a situation in which the program does not give the correct result for some interleavings of the operations of multiple goroutines
    - Race conditions are pernicious because they may remain **latent** in a program and appear infrequently, perhaps only under heavy load or when certain compilers, platforms, or architectures. This makes them hard to reproduce and diagnose
- A data race occurs whenever 2 goroutines access the same variable concurrently and at least of one the accesses is a **write**

    ```go
    package bank
    var balance int
    // 2 operation
    // read: balance + amount
    // write: balance = ...
    func Deposit(amount int) { balance = balance + amount}
    func Balance() int { return balance }

    go func() {
        bank.Deposit(200)                       // A1
        fmt.Println("=", balance.Balance())     // A2
    }()
    go bank.Deposit(100)                        // B
    ```

    - Since the steps `A1` and `A2` occur concurrently with `B`, we cannot predict the order in which they happen. **Intuitively**, it might seen that there're only 3 possible orderings, which we'll call "Alice first", "Bob first", and "Alice/Bob/Alice". In all cases the final balance is $300. The only variation is whether Alice's balance slip includes Bob's transaction or not, but the customers are satisfied either way
    - There's a **fourth** possible outcome, in which Bob's deposit occurs in the middle of Alice's deposit, after the balance has been read (`balance + amount`) but before it has been updated (`balance = ...`), causing Bob's transaction to disappear. This is because Alice's deposit operation `A1` is really a sequence of 2 operations, a read and a write; call them `A1r` and `A1w`
        - `A1r` 读到 balance 值为 0；`B` 完成存款操作；`A` 中 `balance` 仍为 `A1r` 读到的 0，`A1w` 执行 `balance = 0 + 200`；最终 `balance` 是 200，B 的存款操作丢失
    - This program contains a particular kind of race condition called *data race*. ***A data race occurs whenever 2 goroutines access the same variable concurrently and at least one of the accesses is a write***
- Things get messier if the data race involves a variable of a type that is larger than a single machine word

    ```go
    var x []int
    go func() { x = make([]int, 10) }()
    go func() { x = make([]int, 1000000) }()
    x[999999] = 1
    ```

    - `x` could be nil, or a slice of length 10, or a slice of length 1000000. If the **pointer** comes from the first call to `make` and the length comes from the second, `x` would be a slice whose nominal length is 1000000 but whose underlying array has only 10 elements. In this eventuality, storing a element 999999 would clobber an arbitrary faraway location, with consequences that are impossible to predict and hard to debug and localize. The semantic minefield is called *undefined behavior* and is well known to C programmers; fortunately it's rarely as troublesome in Go as in C
        - There're 3 parts to a slice: the pointer, the length, and the capacity
- Even the notion that a concurrent program is an interleaving of several sequential programs is a false intuition. Data races may have even stranger outcomes
    - Many programmers will occasionally offer justifications for known data races in their programs: "the cost of mutual exclusion is too high", "this logic is only for logging", "I don't mind if I drop some messages", and so on. The absence of problems on a given compiler and platform may give them false confidence
    - A good rule of thumb is *there's no such thing as a benign data race*
- 3 ways to avoid a data race
    1. Not to write variable

        ```go
        // Lazily populated as each is requested for the first time

        var icons = make(map[string]image.Image)
        func loadIcon(name string) image.Image
        // NOTE: not concurrency-safe!
        func Icon(name string) image.Image {
            icon, ok := icons[name]
            if !ok {
                icon = loadIcon(name)
                icons[name] = icon
            }
            return icon
        }

        var icons = map[string]image.Image{
            "spades.png": loadIcon("spades.png"),
            "hearts.png": loadIcon("hearts.png"),
            "diamonds.png": loadIcon("diamonds.png"),
            "clubs.png": loadIcon("clubs.png"),
        }
        // Concurrency-safe.
        func Icon(name string) image.Image { return icons[name] }
        ```

        - In the concurrency-safe one, `icons` variable is **assigned during package initialization**, which *happens before* the program's `main` function starts running. Once initialized, `icons` is never modified
        - Cannot use this approach if updates are essential
    2. Avoid accessing the variable from multiple goroutines (variable *confined* to a single goroutine, `ch9/bank1`)
        - Since other goroutines cannot access the variable directly, they must **use a channel** to send the confining goroutine a request to query or update the variable. This is what's meant by the Go mantra **"Do not communicate by sharing memory; instead, share memory by communicating"**
        - A goroutine that brokers access to a confined variable using channel requests is called a *monitor goroutine* for that variable
        - Even when a variable cannot be confined to a single goroutine for its entire lifetime, confinement may still be a solution to the problem of concurrent access
            - For example, it's common to share a variable between goroutines in a pipeline by passing its address from one stage to the next over a channel. If each stage of the pipeline refrains from accessing the variable after sending it to the next stage, then all accesses to the variable are sequential. In effect, the variable is confined to one stage of the pipeline, then confined to the next, and so on. This discipline is sometimes called *serial confinement*
                
            ```go
            type Cake struct{ state string }
            func baker(cooked chan<- *Cake) {
                for {
                    cake := new(Cake)
                    cake.state = "cooked"
                    cooked <- cake // baker never touches this cake again
                }
            }
            func icer(iced chan<- *Cake, cooked <-chan *Cake) {
                for cake := range cooked {
                    cake.state = "iced"
                    iced <- cake // icer never touches this cake again
                }
            }
            ```
        
    3. Allow may goroutines to access the variable, but only one at a time (*mutual exclusion*)
# Mutual exclusion: `sync.Mutex`
- Use a channel of capacity of 1 to ensure that at most 1 goroutine accesses a shared variable at a time. A semaphore that counts only to 1 is called a *binary semaphore* (`ch9/bank2`)
    - This pattern of mutual exclusion is so useful that it's supported directly by the `Mutex` type from `sync` package. Its `Lock` method acquires the token (called a lock) and its `Unlock` method releases it (`ch9/bank3`)
        - Each time a goroutine accesses the shared variables, it must call the mutex's `Lock` method to **acquire a exclusive lock**. If some other goroutine has acquire the lock, this operation will **block** util the other goroutine calls `Unlock` and the lock becomes available again. The mutex *guards* the shared variables. By convention, the variables guarded by a mutex are declared immediately after the declaration of the mutex itself. If you deviate from this, be sure to document it
    - The region of code between `Lock` and `Unlock` in which a goroutine is free to read and modify the shared variables is called a *critical section*. The lock holder's call to `Unlock` happens before any other goroutine can acquire the lock for itself
    - It's essential that the goroutine release the lock **once** it is finish, **on all paths** through the function, including error path
        - In more complex critical sections, especially those in which errors must be dealt with by returning early, it can be hard to tell that calls to `Lock` and `Unlock` are strictly paired on all paths. By deferring a call to `Unlock`, the critical section **implicitly extends to the end of the current function**, freeing us from having to remember to insert `Unlock` calls in one or more places far from the the call to `Lock`. A deferred `Unlock` will run even if the critical section panics, which may be important in programs that make use of `recover`
        - A `defer` is marginally more expensive than an explicit call to `Unlock`, but not expensive enough to justify less clear code. As with concurrent programs, favor clarity and resist premature optimization. Where possible, use `defer` and let critical sections to extend to the end of a function
- A common concurrency pattern
    - A set of exported functions encapsulates one or more variables so that the only way to access the variables is through these functions. Each function acquires a mutex lock at the beginning and releases it at the end, thereby ensuring that the shared variables are not accessed concurrently. This arrangement of functions, mutex lock, and variables are called a *monitor*
- `Withdraw`

    ```go
    // Note: not atomic
    func Withdraw(amount int) bool {
        Deposit(-amount)
        // gap1
        if Balance() < 0 {
            // gap2
            Deposit(amount)
            return false // insufficient funds
        }
        return true
    }
    ```

    - When an excessive withdrawal is attempted, the balance transiently dips below zero. This may causes a concurrent withdrawal for a modest sum to be spuriously rejected. The problem is that `Withdraw` is not *atomic*: it consists of a sequence of 3 separate operations, each of which acquires and then releases the mutex lock, but nothing locks the whole sequence
    
    ```go
    // Note: incorrect
    func Withdraw(amount int) bool {
        mu.Lock()
        defer mu.Unlock()
        Deposit(-amount)
        if Balance() < 0 {
            Deposit(amount)
            return false // insufficient funds
        }
        return true
    }
    ```

    - Mutex locks are **not re-entrant** - it's not possible to lock a mutex that's already locked - this leads to a deadlock where nothing can proceed, and `Withdraw` blocks forever
        - The purpose of a mutex is to ensure that certain invariants of the shared variables are maintained at critical points during program execution. One of the invariants is "no goroutine is accessing the shared variables", but there may be additional invariants specific to the data structures that the mutex guards.
- A common solution (to achieve atomic) is to divide a function into 2: an unexported function that assumes the `Lock` is already held and does the real work, and an exported function that acquires the lock before calling the unexported one

    ```go
    func Withdraw(amount int) bool {
        mu.Lock()
        defer mu.Unlock()
        deposit(-amount)
        // no gap
        if balance < 0 {
            // no gap
            deposit(amount)
            return false // insufficient funds
        }
        return true
    }
    func Deposit(amount) int {
        mu.Lock()
        defer mu.Unlock()
        deposit(amount)
    }
    func Balance() int {
        mu.Lock()
        defer mu.Unlock()
        return balance
    }
    func deposit(amount int) { balance += amount }
    ```

    - Encapsulation by reducing unexpected interactions in a program helps to maintain data structure invariants. For the same reason, encapsulation also helps to maintain concurrency invariants. **When using a mutex, make sure that both it and the variables it guards are not exported**, whether they are package-level variables or fields of a struct
# Read/Write mutexes: `sync.RWMutex`
- `sync.RWMutex` is a *multiple readers, single writer* lock, a special kind of lock that allows read-only operations to proceed in parallel with each other, but write operations to have fully exclusive access

    ```go
    var mu sync.RWMutex
    var balance int
    func Balance() int {
        mu.RLock() // readers lock
        defer mu.RUnlock()
        return balance
    }
    ```

- `Rlock` can be used only if there are no writes to shared variables in the critical section. In general, we should not assume that logically read-only functions or methods don't also update some variables. A method that appears to be a simple accessor might also increment an internal usage counter, or update a cache so that repeat calls are faster
    - If in doubt, use an exclusive `Lock`
- It's only profitable to use an `RWMutex` when most of the goroutine that acquire the lock are readers, and the lock is under *contention*, that is, **goroutines routinely have to wait to acquire it**. An `RWMutex` requires more complex internal bookkeeping, making it slower than a regular mutex for uncontended locks
# Memory synchronization
- Unlike `Deposit`, `Balance` consists only a single operation, so there's no danger of another goroutine executing "in the middle" of it. The reason `Balance` needs mutual exclusion, either channel-based or mutex-based
    1. It's equally important that `Balance` not execute **in the middle of some other operation** like `Withdraw`
        - `Balance` 不加锁在 `Withdraw` 中执行会读到中间状态时的值
    2. Synchronization is about more than just the order of execution of multiple goroutines; synchronization also affects memory
- In a modern computer, there may be dozens of **processors**, each with its **own local cache of the main memory**. For efficiency, writes to memory are buffered within each processor and flushed out to main memory when necessary. They may even be committed to main memory in a **different order** than they were written by the writing goroutine
- Synchronization primitives like channel communications and mutex operations cause the processor to **flush out and commit all its accumulated writes** so that the effects of goroutine execution up to that point are guaranteed to be visible to goroutines running on other processors

    ```go
    // data race
    var x, y int
    go func() {
        x = 1                   // A1
        fmt.Print("y:", y, " ") // A2
    }()
    go func() {
        y = 1                   // B1
        fmt.Print("x:", x, " ") // B2
    }()

    // to be expected
    // y:0 x:1
    // x:0 y:1
    // x:1 y:1
    // y:1 x:1

    // these 2 outcomes might come as a surprise
    // 两个协程运行与两个 cpu，A1，B1 执行完，但这些 cache 都还来不及 flush out 到 main memory
    // x:0 y:0
    // y:0 x:0
    ```

    - Depending on the compiler, CPU, and many other factors, they can happen
    - **Within a single goroutine**, the effects of each statement are guarantee to occur in the order of execution; goroutines are *sequentially consistent*. But in the absence of explicit synchronization using a channel or mutex, there's **no guarantee that events are seen in the same order by all goroutines**. Although goroutine A must observe the effect of the write `x=1` before it reads the value of `y`, it does not necessarily observe the write to `y` done by goroutine B
- It's tempting to try to understand concurrency as if it corresponds to some interleaving of the statements of each goroutine, but this is not how a modern compiler or CPU works
    - Because the assignment and the `Print` refer to different variables, a compiler may conclude that the order of these 2 statements **cannot affect the result**, and **swap** them. If the 2 goroutines execute on different CPUs, each with its own cache, writes by one goroutine are not visible to the other goroutine's `Print` until the caches are synchronized with main memory
- All these concurrency problems can be avoided by the consistent use of simple, established patterns. Where possible, confine variables to a single goroutine; for other variables, use mutual exclusion
# Lazy initialization: `sync.Once`
- It's a good practice to defer an expensive initialization step until the moment it's needed. Initializing a variable up front increases the start-up latency of a program and is unnecessary if execution doesn't always reach the part of the program that uses that variable
- Example

    ```go
    var icon map[string]image.Image
    func loadIcons() {
        icons = map[string]image.Image{
            "spades.png": loadIcon("spades.png"),
            "hearts.png": loadIcon("hearts.png"),
            "diamonds.png": loadIcon("diamonds.png"),
            "clubs.png": loadIcon("clubs.png"),
        }
    }
    // Note: not concurrency-safe!
    func Icon(name string) image.Image {
        if icons == nil {
            loadIcons()
        }
        return icons[name]
    }
    ``` 

    - `Icon` consists of multiple steps: it tests whether `icons` is nil, then it loads the icons, then it updates `icon` to a non-nil value. Intuition might suggest that the worst possible outcome of the race condition is that `loadIcons` function is called several times. The intuition is wrong. In the absence of explicit synchronization, the compiler and CPU are free to **reorder** accesses to memory in any number of ways, so long as the behavior of each goroutine itself is sequentially consistent
    - One possible reorder

        ```go
        var icon map[string]image.Image
        func loadIcons() {
            icons = make(map[string]image.Image)
            icons["spades.png"] = loadIcon("spades.png")
            icons["hearts.png"] = loadIcon("hearts.png")
            icons["diamonds.png"] = loadIcon("diamonds.png")
            icons["clubs.png"] = loadIcon("clubs.png")
        }
        ```

        - It stores the empty map in the `icons` variable before populating it. Consequently, a goroutine finding `icons` to be non-nil may not assume that the initialization of the variable is complete (使得的原先 `icons` 判空的逻辑不对)
    - The simplest correct to ensure that all goroutines observe the effects of `loadIcons` is to synchronize them using a mutex

        ```go
        var mu sync.Mutex // guards icons
        var icons map[string]image.Image
        func Icon(name string) image.Image {
            mu.Lock()
            defer mu.Unlock()
            if (icons == nil) {
                loadIcons()
            }
            return icons[name]
        }
        ```

        - The cost of enforcing mutually exclusive access to `icons` is that 2 goroutines cannot access the variable concurrently, even once the variable has been safely initialized and will never be modified again
    - Using a multiple-readers lock

        ```go
        var mu sync.RWMutex // guards icons
        var icons map[string]image.Image
        func Icon(name string) image.Image {
            // 初始化完成后支持多个协程同时读
            mu.RLock()
            if icons != nil {
                icon := icons[name]
                mu.RUnlock()
                return icon
            }
            mu.RUnlock()

            // lock 写 icons
            mu.Lock()
            if icons == nil { // must recheck for nil
                loadIcons()
            }
            icon := icons[name]
            mu.Unlock()
            return icon
        }
        ```

        - There're 2 critical sections. The goroutine first acquires a reader lock, consults the map, then releases the lock. In an entry was found (the common case), it is returned. If no entry was found, the goroutine acquires a writer lock. There's no way to upgrade a shared lock to an exclusive one without first releasing the shared lock, so we must **recheck** the `icons` in case another goroutine already initialized it in the interim
        - This pattern gives us greater concurrency but is complex and thus error-prone
- The `sync` package provides a specialized solution to the problem of one-time initialization: `sync.Once`. Conceptually, a `Once` consists of a mutex and a boolean variable variable that records whether initialization has taken place; the mutex guards both the boolean and the client's data structures. The sole method, `Do`, accepts the initialization function as its argument

    ```go
    var loadIconOnce sync.Once // guards icons
    var icons map[string]image.Image
    func Icon(name string) image.Image {
        loadIconOnce.Do(loadIcons)
        return icons[names]
    }
    ```

    - Each call to `Do(loadIcons)` locks the mutex and check the boolean value. In the first call in which the variable is false, `Do` calls `loadIcons` and sets the variables to true. Subsequent calls do nothing, but mutex synchronization ensures that the effects of `loadIcons` on memory (specially, `icons`) become visible to all goroutine
- Using `sync.Once` this way, we **can avoid sharing variables with other goroutines until they have been properly constructed**
# The race detector
- Even with the greatest of care, it's all too easy to make concurrency mistakes. Go runtime and toolchain are equipped with a sophisticated and easy-to-use dynamic analysis tool - the *race detector*
- Add `-race` flag to `go build`, `go run`, or `go test`. This causes the compiler to build a modified version of the application or test with additional instrumentation that effectively records all accesses to shared variables that  occurred execution, along with the identity of the goroutine that read or wrote the variable
    - In addition, the modified program records all synchronization events, such as `go` statements, channel operations, and calls to `(*sync.Mutex).Lock`, `(*sync.WaitGrout).Wait`, and so on
    - > The complete set of synchronization events is specified by the *The Go Memory Model* document that accompanies the language specification
- The race detector studies this stream of events, looking for cases in which one goroutine reads or writes a shared variable that was most recently written by a different goroutine **without an intervening synchronization operation**. This indicates a concurrent access to the shared variable, and thus a data race. The tool prints a report that includes the identity of the variable, and the stacks of active function calls in the reading goroutine and the writing goroutine
- The race detector reports all data races that were actually executed. However, it can only detect race conditions that occur during a run; it cannot prove that none will ever occur. For best results, make sure that your tests exercise your packages using concurrency
- Due to extra bookkeeping, a program built with race detection needs more time and memory to run, but the overhead is tolerable even for many production jobs. For infrequently occurring race conditions, letting the race detector do its job can save hours or days or debugging
# Example: concurrent non-blocking cache
- This is the problem of *memorizing* a function, that is, caching the result of a function so that it need be computed only once. The solution will be concurrency-safe and will avoid the contention associated with designs based on a single lock for the whole cache
- Use the `httpGetBody` function as an example of the type of function we might want to memorize. It makes an HTTP GET request and reads the request body. Calls to this function are relatively expensive, so we'd like to avoid repeating them unnecessarily

    ```go
    func httpGetBody(url string) (interface {}, error) {
        resp, err := http.Get(url)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        return ioutil.ReadAll(resp.Body)
    }
    ```

- `ch9/memo1`
    - The reference to `memo.go:30` tells us 2 goroutines have update the `cache` map without any intervening synchronization
- `ch9/memo2`
    - The simplest way to make the cache concurrency-safe is to use monitor-based synchronization. Add a mutex to the `Memo`, acquire the mutex lock at the start of `Get`, and release it before `Get` returns, so that the 2 `cache` operations occur within the critical section
    - By holding the lock for the duration of each call to `f`, `Get` serializes all the I/O operations we intended to parallelize. What we need is a non-blocking cache, one that does not serialize calls to the function it memorizes
- `ch9/memo3`
    - Acquire the lock twice: once for the lookup, and then a second time for the update is the lookup returned nothing. In between, other goroutines are free to use the cache
    - The performance improves again. But some URLs may be fetched twice. This happens when 2 or more goroutines call GET for the same URL at **about** the same time. Both consult the cache, find no value, and then call the slow function `f` (因为另一个稍早发出的请求的结果也没返回). Then both of them update the map with the result they obtained. One of the results is overwritten by the other
    - Ideally we'd like to avoid this redundant work. This feature is sometimes called *duplicate suppression*
- `ch9/memo4`
    - Each map element is a pointer to an `entry` struct. Each `entry` contains the memorized result of a call to the function `f`, as before, but it additionally contains a channel called `ready`. Just after the `entry`'s result has been set, this channel will be closed, to broadcast to any other goroutines that it's now safe for them to read the results from the `entry`
- `ch9/memo5`
    - An alternative implemetation in which the map variable is confined to a monitor goroutine to which calleers of `Get` must send a message
    - The `Memo` type consists of a channel, `requests`, through which the caller of `Get` communicates with the monitor goroutine. The element type of the channel is a `request`. Using this structure, the caller of `Get` sends the monitor goroutine both the key, that is, the argument to the memoized function, and another channel, `response`, over which the result should be sent back when it becomes available. This channel will carry only a single value
    - The `Get` method creates a response channel, puts it in the request, sends it to the monitor goroutine then immediately receives from it
    - The `cache` variable is confined to the monitor goroutine `(*Meno).server`
    - In a similar manner to the mutex-based version, the first request for a given key becomes responsible for calling the function `f` on that key, storing the result in the `entry`, and broadcasting the readiness of the `entry` by closing the `ready` channel. A subsequent request for the same key finds the existing `entry` in the map, waits for the result to become ready, and sends the result through the response channel to the channel goroutine that called `Get`
    - The `call` and `deliver` methods must be called in their own goroutines to ensure that the monitor goroutine does not stop processing new requests
- It's possible to build many concurrent structures using either of the 2 approaches - shared variables and locks, or communicating sequential process - without excessive complexity
    - It's not always obvious which approach is preferable in a given situation, but it's worth knowing how they correspond. Sometimes switching from one approach to the other can make the code simpler
# Goroutines and threads
- Although the differences between goroutines and operating system threads are essentially quantitative, a big enough quantitative difference becomes a qualitative one 
## Growable stacks
- Each OS thread has a fixed-size block of memory (often as large as **2MB**) for its *stack*, the work area where it saves the local variables of function calls that are in progress or temporarily suspended while another function is called
    - This fixed-size stack is simultaneously too much and too little. A 2MB stack would be a huge waste of memory for a little goroutine, such as one that merely waits for a `WaitGroup` then closes a channel. It's not uncommon for a Go program to create hundreds of thousands o goroutines at one time, which would be impossible with stacks this large
- Yet despite their size, fixed-size stacks are not always big enough for the most complex and deeply recursive of functions. Changing the fixed size can improve space efficiency and allow more threads to be created, or it can enable more deeply recursive functions, but it cannot do both
- A goroutine starts life with a small stack, typically **2KB**. A goroutine's stack holds the local variables of active and suspended function calls, is not fixed; it grows and shrinks as needed. The size limit for a goroutine stack may be as much as **1GB**
## Goroutine scheduling
- OS threads are scheduled by the OS kernel. Every few milliseconds, a **hardware timer** interrupts the processor, which causes a kernel function called *`scheduler`* to be invoked
    - This function suspends the currently executing thread and saves its registers in memory, looks over the list of threads and decides which one should run next, restores that thread's registers from memory, then resumes the execution of that thread
- Because OS threads are scheduled by the kernel, passing control from one thread to another requires a full *context switch*, that is, saving the state of one user thread to memory, restoring the state of another of another, and updating the scheduler's data structures
    - This operation is slow, due to its **poor locality** and the number of **memory accesses** required, and has historically only gotten worse as the number of CPU cycles required to access memory has increased
- The Go runtime contains its own scheduler that uses a technique known as ***`m:n` scheduling***, because it multiplexes (or schedules) `m` goroutines on `n` OS threads
    - The job of the Go scheduler is analogous to that of the kernel scheduler, but it is concerned only with the **goroutines of a single Go program**
- The Go scheduler is not invoked periodically by a hardware timer, but implicitly by certain Go language constructs. For example, when a goroutine calls `time.Sleep` or blocks in a channel or mutex operation, the scheduler puts it to sleep and runs another goroutine until it's time to wake the first one up
    - Because it doesn't need a switch to kernel context, rescheduling a goroutine is much cheaper that rescheduling a thread
## `GOMAXPROCS`
- The Go scheduler uses a parameter called `GOMAXPROCS` to determine how many OS **threads** may be actively executing Go code simultaneously. Its default value is the number of CPUs on the machine. It's the `n` in `m:n scheduling`
    - On a machine with 8 CPUs, the scheduler will schedule Go code on up to 8 OS threads at once
- Goroutines that are sleeping or blocked in a communication do not need a thread at all. Goroutines that are blocked in I/O or other system calls or are calling non-Go functions, do need an OS thread, but `GOMAXPROCS` need not account for them
- Can explicitly control this parameter using the `GOMAXPROCS` environment variable or the `runtime.GOMAXPROCS` function
- The effect of `GOMAXPROCS`

    ```go
    for {
        go fmt.Print(0)
        fmt.Print(1)
    }
    // GOMAXPROCS=1 go run main.go
    // GOMAXPROCS=2 go run main.go
    ```

    - In the first run, at most one goroutine was executed at a time. Initially, it was the main goroutine, which prints ones. After a period of time, the Go scheduler put it to sleep and woke up the goroutine that prints zeros, giving it a turn to run on the OS thread
    - In the second run, there were 2 OS threads available, so both goroutines ran simultaneously, printing digits at about the same rate
    - > Many factors are involved in goroutine scheduling, and the runtime is constantly evolving, so the results may differ from the ones in the book
## Goroutines have no identity
- In most operating systems and programming language that support multithreading, the current thread has a distinct identity that can be easily obtained as an ordinary value, typically an integer or pointer
    - This makes it easy to build an abstraction called *thread-local storage*, which is essentially a global map keyed by the thread identity, so that each thread can store and retrieve values independent of other threads
- Goroutines has no notion of identity that is accessible to the programmer. This is by design, since thread-local storage tends to be abused
    - For example, in a web server implemented in a language with thread-local storage, it's common for many functions to find information about the HTTP request on whose behalf they are currently working by looking in that storage. However, just as with programs that rely excessively on global variables, this can lead to an unhealthy "action at a distance" in which the behavior of a function is not determined by its argument alone, but by the identity of the thread in which in runs. Consequently, if the identity of the thread should change–some worker threads are enlisted to help, say–the function misbehaves mysteriously
- Go encourages a simpler style of programming in which parameters that affect the behavior of a function are explicit. Not only does this make programs easier to read, but it lets us freely assign subtasks of a given function to many different goroutine without worrying about their identity