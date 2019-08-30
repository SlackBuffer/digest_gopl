package memo

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client want a single result
}

// The caller of Get communicates with the monitor goroutine through requests.
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
func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
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

// 会有多个 Get 协程
func (memo *Memo) Get(key string) (interface{}, error) {
	responseD := make(chan result)
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
