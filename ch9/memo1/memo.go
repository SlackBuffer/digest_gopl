// Package memo provide a concurrency-unsafe memorization of a function of type Func.
package memo

// Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	cache map[string]result
}

// Func is the type of the function to memorize
type Func func(key string) (interface{}, error)

type result struct {
	value interface{} // []byte
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Not concurrency-safe!
func (memo *Memo) Get(key string) (interface{}, error) { // key is url, memo.f is httpGetBody
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
