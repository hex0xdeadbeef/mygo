package chapter9

import "sync"

var (
	rwMu sync.RWMutex
)

func BalanceRWMutex() int {
	rwMu.RLock()
	defer rwMu.RUnlock()

	return balance
}
