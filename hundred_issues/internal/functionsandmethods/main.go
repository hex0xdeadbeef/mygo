package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"unicode"
)

/*
	FUNCTIONS AND METHODS

	DEFINITIONS
1. Function wraps a sequence of instructions into a module that can be invoked in another place. It can accepts some input data and provide a caller with another output data.
2. Method - is a Function that is bound to a specific type. This type is called "receiver" and can be as a poiner as a literal.

	HOW TO CHOOSE A RECEIVER'S TYPE
1. Using a literal receiver, Go creates a copy of a value and passes it to a method. Any changes remain local to this method and the source object remains unchanged.
2. Using a pointer to a receiver, Go creates a copy of a pointer and the internal changes of method will be reflected on the receiver
3. The rules of choosing a type of a receiver:
	1.1 Receiver MUST be a pointer
		1) If a method changes the receiver. This rule is apliable even if the receiver is a slice and the method inserts the elems.
		2) If a receiver has the field(s) that cannot be copied (for example: sync.Xxx type(s))
	1.2 Receiver SHOULD be a pointer
		1) If a receiver is a large object. It prevents a creation of a huge copy of the instance. A method to find out whether an object is big is benchmarking

	2.1 Receiver MUST be a value
		1) If we need to provide immutability of an instance
	2.2 Receiver SHOULD be a value
		1) The receiver is a slice that needn't be changed
		2) The receiver is a tiny array | struct that is a type of value without fields that cannot be changed
		3) The receiver is a basic type (for example: int, float64, string)
4. We must not mix the receivers in the methods. It must be avoided in general.

	AVOID NAMED RESULTING PARAMETER USING
1. When we make the resulting parameters named, they get a form of the usual variables and get the zero values corresponding to theirs types.
2. With the named params of resulting params we can call an empty return operator, without any additions. In this case the current state of the resulting params will be returned.
3. The cases when we need to use resulting named params:
	1) When we cannot guess the purposes of the resulting params in an interface. It can be replaced with creation of a structure with an elaborative fields
	reflecting the purp ose of them.
	2) Some resulting fields have the same type.
	3) Named params should be used when there's a clear benefit
	4) When we need to postprocess an error occured in the defer section to propagate the error to a caller.
4. Empty return operator is addmitted in short functions because it doesn't harm the readability. In the case of long functions it must be avoided. We cannot mix both ways.

	NAMED PARAMS' SIDE EFFECT
1. When we used named params we need to remember to assign all the relevant vals to them.

	RETURNING NILLED RECEIVER
1. A receiver of pointer can be equal to nil.
2. The method can be invoked on a nil pointer.
3. Method in Go - syntax sugar for a function that has a receiver as a first parameter.
4. When we return a nil pointer that is wrapped by an interface type, we return a non-nil value in fact. To beat this problem we only need to return the pointer used
only if it's not nil, if it's not, we return nil directly explicitly pointing that the result is nil value.
5. An interface that is constructed from a nil receiver is not the nil interface value.

	USING A FILENAME AS A PARAMETER OF A FUNCTION
1. To use a filename in a function is not a best practice. It results in difficulties while writing unit-tests.
2. The cons of this approach:
	1) All the unit tests require theirs own file.
	2) This function isn't reusable.
3. To beat this problem we should wrap this logic into an abstraction. The pros of this approach:
	1) The function abstractes from the source of data
	2) We can substitue args with the different objects

	IGNORING ARGUMENTS EVALUATION IN DEFER STATEMENT
1. The arguments of defer are evaluated at the place when defer statement is arranged.
2. To beat the in-place evaluation of defer args we can use the following ways:
	1) To pass the deferred functions pointers to args, but in this case we narrow the usability of the functions. It requires changes in functions signatures.
	2) Another way to suppres in-place evaluation of args for deferred function is to use closures. In this case the closure gets a fresh value for its execution.

	DEFER VALUE/POINTER RECEIVERS
1. The same logic of in-place execution of deferred funcs' args is applied to using struct funcs.
*/

func main() {
	// AnotherCustomerUsage()

	// NamedParamUsageA()

	// NilReturningUsage()

	// FicticiousTypeUsage()

	// DeferUsageA()

	// DeferUsageFixedA()

	// ClosureUsage()

	// DeferUsageAFixedB()

	// DeferTypeUsageA()

	// DeferTypeUsageB()

}

type customer struct {
	balance float64
	mu      *sync.Mutex
}

// c argument in this method will be a copy of an argument passed in
func (c customer) addLiteral(v float64) customer {
	c.balance += v

	// returns an updated customer val
	return c
}

// c in this method will be a copy of a pointer passed in, the internal changes in this method will be reflected on the receiver
func (c *customer) addPointer(v float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.balance += v
}

type (
	data struct {
		balance float64
	}

	anotherCustomer struct {
		d  *data
		mu *sync.Mutex
	}
)

// Since the field d is defined as a pointer to a data, the source balance will be changed
func (ac anotherCustomer) add(v float64) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.d.balance += v
}

func AnotherCustomerUsage() {
	ac := anotherCustomer{d: &data{balance: 100}, mu: &sync.Mutex{}}

	ac.add(50)
	fmt.Println(ac.d.balance)
}

func NamedParamUsageA() (b int) {
	switch rand.Intn(2) {
	case 1:
		b = 100
	}

	// 100 or 0
	fmt.Println(b)
	return
}

type locator interface {
	// In this method we cannot guess the purpose of the two first resulting arguments, so we must to name them
	// getCoordinates(address string) (float32, float32, error)
	getCoordinates(add string) (latitude, longitude float32, err error)
}

func ReadFull(r io.Reader, buf []byte) (n int, err error) {
	for len(buf) > 0 && err == nil {
		var (
			nr int
		)

		nr, err = r.Read(buf)
		n += nr
		buf = buf[nr:]
	}

	return
}

type locatorContext interface {
	getCoordinates(ctx context.Context, add string) (longitude, latitude float64, err error)
}

type loc struct {
}

func (l *loc) getCoordinates(ctx context.Context, add string) (longitude, latitude float64, err error) {

	if !l.isValidAddr(add) {
		return 0, 0, errors.New("invalid address")
	}

	// We cannot write this code because we ignore an error of ctx. We must assign it to the err named variable.
	// if err = ctx.Err(); err != nil {
	// 	return 0, 0, err
	// }

	// To beat the problem we can define and assign to this defined err variable an error returned. In this case the err shades the named param err
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}

	// Or use the named parameter parameter
	if err = ctx.Err(); err != nil {
		return 0, 0, err
	}

	// Or use the predefined named parameter with the empty operator. But in this case we mix both forms of returning values during usage named params.
	if err = ctx.Err(); err != nil {
		return
	}

	// getting coordinates and returning them ...
	return 100.0, 120.0, nil
}

func (l *loc) isValidAddr(addr string) bool {
	for _, r := range addr {
		if unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

type MultiError struct {
	errs []error
}

func (m *MultiError) Add(err error) {
	m.errs = append(m.errs, err)
}

func (m *MultiError) Error() string {
	var (
		errsSize int
		b        *strings.Builder
	)

	for _, err := range m.errs {
		errsSize += len(err.Error())
	}
	b.Grow(errsSize)

	for _, err := range m.errs {
		b.WriteString(err.Error())
		b.WriteByte(';')
	}

	return b.String()
}

func NewMultiError() *MultiError {
	return &MultiError{errs: make([]error, 0, 1<<4)}
}

type Customer struct {
	Name string
	Age  int
}

func (c *Customer) Validate() error {
	var (
		m *MultiError
	)

	if c.Age < 0 {
		m.Add(errors.New("age is negative"))
	}

	if len(c.Name) == 0 {
		if m == nil {
			m = &MultiError{}
		}
		m.Add(errors.New("name is nil"))
	}

	// Before returning the value m is wrapped by the error interface
	// The value will be returned by a function is not nil error because the type descriptor of an error is defined, but the value is not
	// return m

	// To beat this problem we need to return m only if it's not nil
	if m != nil {
		return m
	}

	// If it's no, we return nil directly
	return nil
}

func NilReturningUsage() {
	customer := Customer{Age: 33, Name: "John"}
	if err := customer.Validate(); err != nil {
		log.Fatalf("customer is invalid: %#v", err)
	}
}

type Foo struct{}

// In fact this is:
//
//	func Bar(foo *Foo) string {
//		return "bar"
//	}
func (foo *Foo) Bar() string {
	return "bar"
}

func FicticiousTypeUsage() {
	var (
		foo *Foo
	)

	// In fact the function is invoked here. foo.Bar(nil)
	fmt.Println(foo.Bar())
}

func countFileEmptyLinesWrong(fileName string) (n int, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer func() {
		fileClosingErr := file.Close()
		if fileClosingErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; closing file: %w", err, fileClosingErr)
		}

		err = fileClosingErr
	}()

	var (
		scanner = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		// processing
	}

	return n, err
}

// Using an abstraction over reader this function can accept any readers
func countFileEmptyLinesRightWithAbstraction(r io.Reader) (n int, err error) {
	if r == nil {
		return 0, errors.New("an empty reader passed")
	}

	var (
		// bufio.NewScanner(r io.Reader)
		s = bufio.NewScanner(r)
	)

	for s.Scan() {
		// processing
	}

	return n, err
}

// There will not any outputs in deferred statements because the args of defer statement is captured
// at the place when defer is arranged
func DeferUsageA() error {
	const (
		StatusSuccess  = "success"
		StatusErrorFoo = "error_foo"
		StatusErrorBar = "error_bar"
	)
	var (
		status string
	)
	defer notifyA(status)
	defer incrementCounterA(status)

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

// We can suppres in-place evaluation of deferred funcs with passing pointers to target functions
func DeferUsageFixedA() error {
	const (
		StatusSuccess  = "success"
		StatusErrorFoo = "error_foo"
		StatusErrorBar = "error_bar"
	)
	var (
		status string
	)

	defer notifyB(&status)
	defer incrementCounterB(&status)

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

// Another way is using closures
func DeferUsageFixedB() error {
	const (
		StatusSuccess  = "success"
		StatusErrorFoo = "error_foo"
		StatusErrorBar = "error_bar"
	)

	var (
		status string
	)

	defer func() {
		incrementCounterA(status)
		notifyA(status)
	}()

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

func foo() error {
	switch rand.Intn(2) {
	case 1:
		return errors.New("error_foo")
	default:
		return nil
	}
}

func bar() error {
	switch rand.Intn(2) {
	case 1:
		return errors.New("error_bar")
	default:
		return nil
	}
}

func notifyA(status string) {
	fmt.Printf("notify called: %q\n", status)
}

func incrementCounterA(status string) {
	fmt.Printf("incrementCounter called: %q\n", status)
}

func notifyB(status *string) {
	fmt.Printf("notify called: %q\n", *status)
}

func incrementCounterB(status *string) {
	fmt.Printf("incrementCounter called: %q\n", *status)
}

func ClosureUsage() {
	var (
		i int
		j int
	)

	defer func(i int) {
		// j refers to the outer variable j
		// its value evaluated only after invocation of this closure
		fmt.Println(i, j)
	}(i)

	i++
	j++
}

type Worker struct {
	id        string
	telNumber *TelNumber
}

type TelNumber struct {
	number int
}

func (w Worker) PrintValue() {
	fmt.Println(w.id, w.telNumber)
}

func (w *Worker) PrinPointer() {
	fmt.Println(w.id, w.telNumber)
}

func DeferTypeUsageA() {
	worker := Worker{id: "QWERTY", telNumber: &TelNumber{number: 333}}
	defer worker.PrintValue()

	worker.id = "ABCDEFG"
	worker.telNumber.number = 777

	fmt.Println("Function return:", worker.id, worker.telNumber)
	fmt.Print("Defer invocation: ")
}

func DeferTypeUsageB() {
	worker := Worker{id: "QWERTY", telNumber: &TelNumber{number: 333}}
	defer worker.PrinPointer()

	worker.id = "ABCDEFG"
	worker.telNumber.number = 777
	fmt.Println("Function return:", worker.id, worker.telNumber)
	fmt.Print("Defer invocation: ")
}

// It'll rerurn 2 because of calling defer function on the explicit x value
func test() (x int) {
	defer func() {
		x++
	}()

	defer func() {
		x++
	}()

	defer func() {
		x = 10
	}()

	x = 1
	return
}

// It'll return 1 because of explicit return of the value
func anotherTest() int {
	var x int
	defer func() {
		x++
	}()

	defer func() {
		x++
	}()

	x = 1
	return x
}
