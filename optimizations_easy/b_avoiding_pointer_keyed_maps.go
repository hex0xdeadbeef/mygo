package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/*
	AVOIDING LARGE MAPS THAT STORES VALUES WITH POINTER KEYS
1. At the time of scanning trash in the heap the environment is scanning the objects with the pointers and tracking them. If we have a large map (for example: map[string]int), the
	GC should check every string. It happens at each iteration of the GC collection since strings includes pointers.
2. After execution of the function testGCTime() we see the following output:
	GC Took: 109.360708ms !!!
	GC Took: 698.666µs
	GC Took: 1.06675ms !!!
	GC Took: 860.292µs
	GC Took: 721.291µs
	GC Took: 557.625µs
	GC Took: 524.167µs
3. To beat this problem we should replace pointer-keys ubiquitously switching them to a digit-type to dump the workload from the GC. The result is:
	GC Took: 172.375µs
	GC Took: 668.75µs
	GC Took: 562.958µs
	GC Took: 438.875µs
	GC Took: 418.292µs
	GC Took: 295.667µs
4. Compairing two types of map keys we see that the difference is colossal.
5. Using the map in production we need to hash the strings into the integer numbers. The subtlety is the following trick:
	If we allocate a giant []byte, []int and etc. or an array of pointerless structs, Go treats if as "non-scan" so the GC's scanning phase gets a break. It can require gimmlicks
	like using numeric IDs/Offsets instead of pointers, much as if we were implementing an on-disk database, but we can continue to use language constructs to allocate instead of
	having to go to the OS.
*/

const (
	tokenSize = 1 << 8
)

var (
	tokenSymbolsSet []byte
)

func init() {
	const (
		digitsNumber   = 10
		digitStartByte = '0'

		lettersNumber      = 26
		lowercaseStartByte = 'a'
		uppercaseStartByte = 'A'
	)

	tokenSymbolsSet = make([]byte, 0, digitsNumber+lettersNumber*2)

	for i := 0; i < digitsNumber; i++ {
		tokenSymbolsSet = append(tokenSymbolsSet, digitStartByte+byte(i))
	}

	for i := 0; i < lettersNumber; i++ {
		tokenSymbolsSet = append(tokenSymbolsSet, lowercaseStartByte+byte(i), uppercaseStartByte+byte(i))
	}
}

type TokenCache struct {
	body map[Token]struct{}
	mu   *sync.RWMutex
}

func NewTokenCache() *TokenCache {
	return &TokenCache{
		body: initialTokens(),
		mu:   &sync.RWMutex{},
	}
}

func initialTokens() map[Token]struct{} {
	const (
		startCacheBodySize = 1 << 8
	)

	var (
		body = make(map[Token]struct{}, startCacheBodySize)
	)

	for i := 0; i < startCacheBodySize; i++ {
		putUniqueToken(body)
	}

	return body
}

func (tc *TokenCache) Put() {
	tc.mu.Lock()
	putUniqueToken(tc.body)
	tc.mu.Unlock()
}

func (tc *TokenCache) GetAll() map[Token]struct{} {
	tc.mu.RLock()
	copiedBody := make(
		map[Token]struct{},
		len(tc.body),
	)

	for k, v := range tc.body {
		copiedBody[k] = v
	}
	tc.mu.RUnlock()

	return copiedBody
}

type Token string

func putUniqueToken(m map[Token]struct{}) {
	var (
		isPresent bool
		curToken  Token
	)

	curToken = getToken()
	_, isPresent = m[curToken]

	for isPresent {
		curToken = getToken()
		_, isPresent = m[curToken]
	}

	m[curToken] = struct{}{}
}

func getToken() Token {
	var (
		buf = make([]byte, 1<<8)
	)

	for i := 0; i < tokenSize; i++ {
		buf[i] = tokenSymbolsSet[rand.Intn(len(tokenSymbolsSet))]
	}

	return Token(string(buf))
}

const (
	maxTokenCacheSize = 1e8
)

func testGCTimeA() {

	tc := NewTokenCache()

	for i := 0; i < maxTokenCacheSize; i++ {
		tc.Put()
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}

}

func testGCTimeB() {
	var (
		m = make(map[uint64]struct{}, maxTokenCacheSize)
	)

	for i := 0; i < maxTokenCacheSize; i++ {
		fillDigitMap(m)
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}

}

func fillDigitMap(m map[uint64]struct{}) {
	var (
		curKey    = rand.Uint64()
		isPresent bool
	)

	curKey = rand.Uint64()
	_, isPresent = m[curKey]

	for isPresent {
		curKey = rand.Uint64()
		_, isPresent = m[curKey]
	}

	m[curKey] = struct{}{}
}

func timeGC() {
	t := time.Now()
	runtime.GC()
	fmt.Printf("GC Took: %s\n", time.Since(t))
}
