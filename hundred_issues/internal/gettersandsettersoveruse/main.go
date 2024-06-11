package main

import (
	"math/big"
	"sync"
)

type BankAcc struct {
	balance big.Int

	mu *sync.RWMutex
}

func (ba *BankAcc) Balance() big.Int {
	ba.mu.RLock()
	b := ba.balance
	ba.mu.RUnlock()

	return b
}

func (ba *BankAcc) SetBalance(num big.Int) bool {
	ba.mu.RLock()
	if num.Sign() < 0 {
		return false
	}
	ba.mu.RUnlock()

	ba.mu.Lock()
	ba.balance = num
	ba.mu.Unlock()

	return true
}
