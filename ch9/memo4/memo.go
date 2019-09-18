package memo

import "sync"

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// a memo caches the results of calling a Func
type Memo struct {
	f     Func
	mu    sync.Mutex // in bytes.Buffer, the initial value of the struct is a ready-to-use empty buffer; the zero value of `sync.Mutex` is a ready-to-use unlocked mutex
	cache map[string]*entry
}

// Func is the type of the function to memorize
type Func func(key string) (interface{}, error)

// []byte; err
type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

// key is url; memo.f is httpGetBody
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing the value and broadcasting the ready condition.
		// 将 e 置为非 nil，做到 duplicate suppression；但此时还不能读取值，需等广播
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e

		// 开始网络请求前就将锁释放
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		// 此时网络请求的实际数据可能还未返回，但 e 已不是 nil
		memo.mu.Unlock()
		// This operation blocks until the channel is closed.
		<-e.ready // wait for ready condition
	}
	// e.res.value and e.res.err are shared among multiple goroutines. Despite being accessed by multiple goroutines, no mutex lock is necessary.
	// The closing of the ready channel happens before any other goroutine receives the broadcast event,
	// so the write to those variables in the first goroutine happens before they are read by subsequent goroutines. There's no data race
	return e.res.value, e.res.err
}
