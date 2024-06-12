package main

import (
	"fmt"
	"net/http"
)

/*
	FUNCTIONAL OPTIONS IN USE

1. We can use a structure to arrange some optinal params. The mandatory params will be placed as the usual ones.
2. We can use a "Builder" pattern to create configs
*/

func main() {

}

// This way isn't flexible if we wanna make decision based on port argument
func NewServerA(add string, port int) (*http.Server, error) {
	var (
		srvr *http.Server
	)
	// some logic

	return srvr, nil
}

// The ConfigA consists of the optional params for the function.
type ConfigA struct {
	Port int
}

// This way is pleasant 'cause we can decide what to do refering to cfg Port field. Presence of Config param results in more flexible desing of function.
func NewServerB(add string, cfg ConfigA) (*http.Server, error) {
	var (
		srvr *http.Server
	)
	// some logic

	return srvr, nil
}

// But in the case when we need to explicitly point the zero val of port we fails.
// The first way is to make the Port field a pointer to an int val
func EqualityOfZeroedFuncs() {
	c1 := ConfigA{} // Zeroed Port val
	c2 := ConfigA{Port: 0}

	fmt.Println(c1 == c2)
}

// After creation a zeroed val of ConfigB instance we can differentiate the nil val of Port and 0 val of this field.
// This way lacks the flexibility of instance creation.
// 1) The client should create an integer val ad pass into creation literal firstly
// 2)
type ConfigB struct {
	Port *int
}

// This way is pleasant 'cause we can decide what to do refering to cfg Port field. Presence of Config param results in more flexible desing of function.
func NewServerC(add string, cfg ConfigB) (*http.Server, error) {
	var (
		srvr *http.Server
	)
	// some logic

	return srvr, nil
}

func LacksConfigPointerField() {
	// Creation of the integer and the explictit passing it into a cfg struct
	port := 8090
	cfg := ConfigB{&port}

	fmt.Println(cfg)

	// The basic config creation
	srvr, err := NewServerC("localhost:", ConfigB{})
	if err != nil {
		panic(err)
	}
	fmt.Println(srvr)
}
