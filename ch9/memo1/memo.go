// provide a concurrency-unsafe memorization of a function of type Func
package memo

// a memo caches the results of calling a Func
type Memo struct {
	f     Func
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

// not concurrency-safe; multiple goroutines may update the `cache` map
// without any intervening synchronization
// `key` is url; `memo.f` is `httpGetBody`
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
