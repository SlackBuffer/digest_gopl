package bank

var (
	sema    = make(chan struct{}, 1) // a binary semaphore guarding balance
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} // accquire token
	balance = balance + amount
	<-sema // release token
}

func Balance() int {
	sema <- struct{}{}
	b := balance
	<-sema
	return b
}
