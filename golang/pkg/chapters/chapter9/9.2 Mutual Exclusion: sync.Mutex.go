package chapter9

import (
	"sync"
)

var (
	sema             = make(chan struct{}, 1)
	balanceSemaphore int
)

func DepositSemaphor(amount int) {
	sema <- struct{}{}
	balanceSemaphore = balanceSemaphore + amount
	<-sema
}

func BalanceSemaphore() int {
	sema <- struct{}{}
	b := balanceSemaphore
	<-sema

	return b
}

var (
	mu           sync.Mutex
	balanceMutex int
)

func DepositMutex(amount int) {
	mu.Lock()
	defer mu.Unlock()

	if amount >= 0 {
		balanceMutex += amount
	}
}

func BalanceMutex() int {
	mu.Lock()
	defer mu.Unlock()

	return balance
}
