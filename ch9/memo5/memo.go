package memo

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client want a single result
}

/* Memo 由 request 类型的 channel 组成。调用 Memo.Get 的协程将 request 发到该 channel。request 包括要 cache 的 key；response 用于发送是 GET 请求的结果。*/
// The caller of Memo.Get communicates with the monitor goroutine through requests.
type Memo struct{ requests chan request }

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// Func is the type of the function to memorize
type Func func(key string) (interface{}, error)

// result of calling a Func
// []byte; err
type result struct {
	value interface{}
	err   error
}

// New returns a memorization of f. Clients must subsequently call Close. (否则 server 会无法结束)
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

// monitor goroutine
// cache 限定在 (*Memo.Server) 协程。
// client 调用 New 后 server 开始运行。
func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	// 循环一直到 *Memo.Close 关闭 memo.requests
	for req := range memo.requests {
		// 询问 cache 中是否存在对应 key
		e := cache[req.key]
		if e == nil {
			// 第一次请求该 key
			e = &entry{ready: make(chan struct{})}
			// 置为非 nil，做到 duplicate suppression
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

// 请求成功后通知 deliver 协程去取
func (e *entry) call(f Func, key string) {
	// evaluate the function
	e.res.value, e.res.err = f(key)
	// broadcast the ready condition
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// wait for the ready condition
	<-e.ready
	// send the result to the client
	response <- e.res
	// fmt.Printf("within deliver: %T\n", response)
}

// 多个 Get 协程并发
func (memo *Memo) Get(key string) (interface{}, error) {
	responseD := make(chan result)
	// Get 将请求内容写到与 Memo.requests 让 server 处理
	// responseD 本身未被转成单向 channel
	memo.requests <- request{key: key, response: responseD}
	// responseD 指向的地址写入了内容，responseD 本身仍是双向 channel
	res := <-responseD

	/* response := make(chan result)
	// response 本身未被转成单向 channel
	memo.requests <- request{key, response}
	// memo.requests <- request{key: key, response: response}
	// response 指向的地址写入了内容，response 本身仍是双向 channel
	res := <-response
	fmt.Printf("%T\n", response) */

	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

/*
- Get 并发请求 url，每个协程往和 server 协程通信的 channel–memo.requests–里写请求
	- 每个请求包括一个 url 值，一个给 server 协程塞请求结果的 result channel
	- 然后立即从 result channel 读值，读到值前一致保持阻塞，result chanel 虚位以待
- cache 里的每条 entry 是每条记录的实际数据和一个 channel，channel 起到两个作用：第一次收到请求时占位，避免重复请求；请求成功后作广播
- cache 限定在 server 的协程
	1. gr1 协程第一次请求某个 url，server 用 entry{ready: make(chan struct{})} （一个用作广播通知的 channel）放入 cache 占位，作用告诉其它请求同一个 url 的协程请求已由发出，值返回了会发广播通知
		- call 拿回请求结果后，存入 entry.result，通过 Close(entry.ready) 广播通知其它协程，此 key 的值已拿回
		- deliver 阻塞，直到收到广播通知值已取回，将 entry.result 的值写入 request.response channel，供 Get 协程去读
	2. 后续的协程都可以从关闭的 ready channel 得知值已取回
- 并发 Get 往 server channel 写数据等 server 协程里的 range 去消费
	- server 协程里的 range 操作保持了并发性，call 和 deliver 都是在各自的协程里进行
- 并发 Get 请求发完后关闭 memo.requests，告知 server 任务结束
*/
