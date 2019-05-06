package memo

import "sync"

type entry struct {
	res   result
	ready chan struct{} // closed when `res` is ready
}

// a memo caches the results of calling a Func
type Memo struct {
	f     Func
	mu    sync.Mutex
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

// `key` is url; `memo.f` is `httpGetBody`
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// first request for this key
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()
		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // broadcast ready condition
	} else {
		// repeat request for this key
		memo.mu.Unlock()
		// this operation blocks until the channel is closed
		<-e.ready // wait for ready condition
	}
	return e.res.value, e.res.err
}
