package main

import (
	cryptorand "crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"
)

/*
	AVOID USING math/rand, USE crypto/rand
1. When we're working on projects that require generating keys, like for encryption or for creating unique identifiers, the quality and security of those keys are really important.

2. Why not math/rand?
	The math/rand package generates `pseudo-random` numbers.

	This means if we know how the numbers are generated (the seed), we can predict the numbers.

	Even if we seed it with the current time (like time.Nanoseconds()), the amount of unpredictability (entropy) is low because there's not a lot of variation in the current time from one execution to the next.

3. Why crypto/rand?
	crypto/rand provides a way to generate numbers that are cryptographically secure. It's designed to be unpredictable, using sources of randomness provided by our OS, which are much harder to predict.

	crypto/rand is suitable for encryption, authentification, and other security-sensitive ops.


	EMPTY SLICE OR, EVEN BETTER, `nil` SLICE
1. When working with slices in Go, we have two approaches to start with what appears to be an empty slice:
	- Using var keyword `var t []int`. This method declares a slice `t` of type `[]int` without initializing it. The slice is considered nil. This means that it doesn't actually point to any underlying array, Its length (len) and capacity (cap) are both 0.
	- Using slice literal `t := []int{}`. Unlike var declaration, this slice is not nil. It's a slice that points to an underlying array, but that array has no elems.

2. So which one is considered idiomatic?
	A `inl` slice doesn't allocate any memory. It's just a pointer to nowhere, while an empty slice (`[]int{}`) actually allocates a small amount of memory to point to an existing, but empty array. In most cases, this difference is negligible, but for high-performance applications, it could be significant.


	The Go community prefers the `nil` slice approach because it's considered more idiomatic to the languages philosophy of simplicity and zero values.

	Of course, exceptions exist. For example, when working with JSON as null, whereas an empty slice (`t := []int{}`) encodes to an empty JSON array.
		- It's also idiomatic to design our code to treat a non-empty slice, an empty slice, and a nil slice similarly.

	If we're familiar with Go, we may know that for rage, len, append, ... work without panic with a nil slice.e


	ERROR MESSAGES SHOULDN'T BE CAPITALIZED OR END WITH PUNCTUATION MARKS.
1. Why lowercase?
	Error messages often get wrapped or combined with other messages. If an error string starts with a capital letter, it can look odd or out of place when it's in the middle of a sentence.

	Starting with a lowercase letter helps it blend more naturally:
		// application startup failed: failed to initialize module: could not open the database

	This means, any text following %w in a formatted error string is intended to be appended at the end of the whole message.

2. Why no punctual?
	It's to ensure that when one message is appended to another, the result looks like a coherent sentence rather than a jumble of phrases.
*/

func MathRandKey() string {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	buf := make([]byte, 16)
	for i := range buf {
		buf[i] = byte(r.Intn(256))
	}

	return fmt.Sprintf("%x", buf)
}

func CryptoRandKey() string {
	buf := make([]byte, 16)

	_, err := cryptorand.Read(buf)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", buf)
}
