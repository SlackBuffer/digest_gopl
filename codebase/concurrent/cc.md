
```go
// 写入 channel 操作放在子协程里，方便
func main() {
	fn := make(chan string)
	go func(ff chan<- string) {
		for _, f := range []string{"fn1", "fn2", "fn3"} {
			ff <- f
		}
		close(ff)
    }(fn)
    
	fmt.Println(pic(fn))
}
func pic(fn <-chan string) string {
	str := make(chan string)
	var wg sync.WaitGroup
	for f := range fn {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			str <- f
		}(f)
    }
    
	go func() {
		wg.Wait()
		close(str)
	}()

	var ss string
	for s := range str {
		ss += s + "\n"
    }
    return ss
}

// 写入 channel 放在 main 协程里
func main() {
	fn := make(chan string)
	// res := go pic(fn) 写法不对，协程间用 channel 传递数据，传统的直接赋值行不通
	res := make(chan string)
	go pic(fn, res)

	for _, f := range []string{"fn1", "fn2", "fn3"} {
		fn <- f
	}
	close(fn)

	// 没有此处的阻塞，会出现 main 协程返回，子协程未执行完就被终止的情况，导致实际上没有完成预期的操作
	fmt.Println(<-res)
}
func pic(fn <-chan string, res chan<- string) {
	str := make(chan string)
	var wg sync.WaitGroup
	for f := range fn {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			str <- f
		}(f)
	}

	go func() {
		wg.Wait()
		close(str)
	}()

	var ss string
	for s := range str {
		ss += s + "\n"
	}
	res <- ss
}
```