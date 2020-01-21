package memo

import "sync"

// Memo caches the results of calling a Func
type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

// Func is the type of the function to memorize
type Func func(key string) (interface{}, error)

// []byte; err
type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// key is url; memo.f is httpGetBody
func (memo *Memo) Get(key string) (interface{}, error) {
	// 不存在 data race，但同时失去了并发的优势，同一时刻只能有一个协程访问 cache
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()

	return res.value, res.err
}
