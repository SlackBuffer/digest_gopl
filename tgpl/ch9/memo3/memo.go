package memo

import "sync"

// a memo caches the results of calling a Func
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
	// lookup cache 的时候加第一次锁
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()

	/** 两个锁的中间其它协程可以访问 cache
	***		部分 url 可能访问两次，因为另一个稍早发出的请求的结果也没返回，导致更新 cache 的同一个键两次。期望能做到 duplicate suppression
	 */

	// 更新 cache 时加第二次锁
	if !ok {
		res.value, res.err = memo.f(key)

		// Between the 2 critical sections, several goroutines may race to compute f(key) and update the map
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	memo.mu.Unlock()
	return res.value, res.err
}
