package main

import (
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkAvgBalanceB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			wg = &sync.WaitGroup{}
			c  = NewCache()
		)
		c.Fill()

		wg.Add(1)
		go func() {
			wg.Done()

			for i := 0; i < 10_000; i++ {
				for k := range c.balances {
					c.AddBalance(k, multiplierBase*rand.Float64())
				}
			}

		}()

		wg.Add(1)
		go func() {
			wg.Done()
			for i := 0; i < 10_000; i++ {
				c.AvgBalanceB()
			}
		}()

		wg.Wait()
	}
}

func BenchmarkAvgBalanceC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			wg = &sync.WaitGroup{}
			c  = NewCache()
		)
		c.Fill()

		wg.Add(1)
		go func() {
			wg.Done()

			for i := 0; i < 10_000; i++ {
				for k := range c.balances {
					c.AddBalance(k, multiplierBase*rand.Float64())
				}
			}

		}()

		wg.Add(1)
		go func() {
			wg.Done()
			for i := 0; i < 10_000; i++ {
				c.AvgBalanceC()
			}
		}()

		wg.Wait()
	}

}
