# Race conditions
- When we cannot confidently say that one event happens before the other, then the events `x` and `y` are concurrent
- A function is concurrency-safe if it continues to work correctly even when called concurrently, that is, from 2 or more goroutines with no additional synchronization
- A type is concurrency-safe if all its accessible methods and operations are concurrency-safe
- We can make a program concurrency-safe without making every concrete type in that program concurrency-safe
    - Concurrency-safe types are the exceptions rather than the rule
    - You should access a variable concurrently only if the documentation for its type says that it's safe
- We avoid concurrent access to most variables either by confining them to a single goroutine or by maintaining a higher-level invariant of mutual exclusion
- Exported package-level functions are generally expected to be concurrency-safe
    - Since package-level variables cannot be confined to a single goroutine, functions that modify them must enforce mutual exclusion
- Many reasons a function might not work when called concurrently, including dead-lock, livelock, and resource starvation
- A race condition is a situation in which the program does not give the correct result for some interleavings of the operations of multiple goroutines
- Race conditions are pernicious because they may remain **latent** in a program and appear infrequently, perhaps only under heavy load or when certain compilers, platforms, or architectures
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
        bank.Deposit(200)
        fmt.Println("=", balance.Balance())
    }()
    go bank.Deposit(100)
    ```

    - If the second deposit occurs in the middle of the first deposit, after the balance has been read but before it has been updated, it'll cause the second transaction to [ ] disappear
- Things get messier if the data race involves a variable of a type that is larger than a single machine word

    ```go
    var x []int
    go func() { x = make([]int, 10) }()
    go func() { x = make([]int, 1000000) }()
    x[999999] = 1
    ```

    - `x` could be nil, or a slice of length, or a slice of length 1000000
    - There're 3 parts to a slice: the pointer, the length, and the capacity
    - If the pointer comes from the first call to `make` and the length comes from the second, `x` would be a slice whose nominal length is 1000000 but whose underlying array has only 10 elements. In this eventuality, storing a element 999999 would clobber an arbitrary faraway location, with consequences that are impossible to predict and hard to debug and localize
    - The semantic minefield is called undefined behavior
- The notion that a concurrent program is an interleaving of several sequential programs is a false intuition
- A good rule of thumb is there's no such thing as a benign data race
- 3 ways to avoid a data race
    1. Not to write variable
    2. Avoid accessing the variable from multiple goroutines (variable confined to a single goroutine)
        - Since other goroutines cannot access the variable directly, they must use a channel to send the confining goroutine a request to query or update the variable
        - This is what's meant by the Go mantra "Do not communicate by sharing memory; instead, share memory by communicating"
        - A goroutine that brokers access to a confined variable using channel requests is called a **monitor goroutine** for that variable
    3. Allow may goroutines to access the variable, but only one at a time (mutual exclusion)
- It's common to share a variable between goroutines in a pipeline by passing its address from one stage to the next over a channel. If each stage of the pipeline refrains from access the variable after sending it to the next stage, then all accesses to the variable are sequential. In effect, the variable is confined to one stage of the pipeline, then confined to the next, and so on
    - This discipline is sometimes called *serial confinement*

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

# Mutual exclusion: `sync.Mutex`
- Use a channel of capacity of 1 to ensure that at most 1 goroutine accesses a shared variable at a time
- A semaphore that counts only to 1 is called a **binary semaphore**
- This pattern of mutual exclusion is so useful that it's supported directly by the `Mutex` type from `sync` package
    - Its `Lock` method acquires the token (called a lock) and its `Unlock` method releases it
    - Each time a goroutine accesses the shared variables, it must call the mutex's `Lock` method to **acquire a exclusive lock**. If some other goroutine has acquire the lock, this operation will **block** util the other goroutine calls `Unlock` and the lock becomes available again
    - The mutex guards the shared variables
    - By convention, the variables guarded by a mutex are declared immediately after the declaration of the mutex itself
        - If you deviate from this, be sure to document it
- The region of code between `Lock` and `Unlock` in which a goroutine is free to read and modify the shared variables is called a critical section
- It's essential that the goroutine release the lock **once** it is finish, **on all paths** through the function, including error path
    - By deferring a call to `Unlock`, the critical section **implicitly extends to the end of the current function**, freeing us from having to remember to insert `Unlock` calls in one or more places far from the the call to `Lock`
    - A deferred `Unlock` will run even if the critical section panics, which may be important in programs that make use of `recover`
- A `defer` is marginally more expensive than an explicit call to `Unlock`, but not expensive enough to justify less clear code
    - As with concurrent programs, favor clarity and resist premature optimization
    - Where possible, use `defer` and let critical sections to extend to the end of a function
- Mutex locks are **not re-entrant** - it's not possible to lock a mutex that's already locked - this leads to a deadlock where nothing can proceed
    - When a goroutine acquires mutex lock, it may assume that the invariants of the shared variables hold
    - While it holds the lock, it may update the shared variables so that the invariants are temporarily violated
    - When it releases the lock, it must guarantee that the invariants hold again
- A common solution (to achieve atomic) is to divide a function into 2: an unexported function that assumes the `Lock` is already held and does the real work, and an exported function that acquires the lock before calling the unexported one

    ```go
    // not atomic
    func Withdraw(amount int) bool {
        Deposit(-amount)
        if Balance() < 0 {
            Deposit(amount)
            return false // insufficient funds
        }
        return true
    }
    // When an excessive withdrawal is attempted, the balance transiently 
    // dips below zero. This may causes a concurrent withdrawal for a modest 
    // sum to be spuriously rejected
    ```

- Encapsulation by reducing unexpected interactions in a program helps to maintain data structure invariants
- For the same reason, encapsulation also helps to maintain concurrency invariants
- When using a mutex, make sure that both it and the variables it guards are not exported, whether they are package-level variables or fields of a struct
# Read/Write mutexes: `sync.RWMutex`