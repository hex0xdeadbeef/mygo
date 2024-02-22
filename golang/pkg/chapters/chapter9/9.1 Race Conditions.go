package chapter9

import (
	"fmt"
	"image"
	"log"
)

var balance int

func Deposit(amount int) {
	balance = balance + amount
}

func Balance() int {
	return balance
}

func Race() {
	var logger log.Logger

	go func() {
		Deposit(200)
		logger.Print("=", Balance())
	}()

	go Deposit(100)
}

func SliceChanging() {
	var x []int
	go func() { x = make([]int, 10) }()
	go func() { x = make([]int, 1000000) }()
	x[99999] = 1
}

/*
Three ways to avoid data races
*/

var icons = make(map[string]image.Image)

func loadIcon(name string) image.Image {

	return nil
}

// UnsafeIcon() is not concurrency-safe
func UnsafeIcon(name string) image.Image {
	icon, ok := icons[name]
	if !ok {
		icon = loadIcon(name)
		icons[name] = icon
	}

	return icon
}

var prefilledIcons = map[string]image.Image{
	"spades.png": loadIcon("spades.png"),
	"hearts.png": loadIcon("hearts.png"),
	// ...
}

// SafeIcon() is concurrency-safe
func SafeIcon(name string) image.Image {
	return icons[name]
}

type withdrawal struct {
	amount  int
	succeed chan bool
}

var (
	// send amount to deposit
	deposits = make(chan int)
	// receive balance
	balances = make(chan int)
	// withdraw funds
	withdrawals = make(chan withdrawal)
)

func init() {
	go teller()
}

func DepositSafe(amount int) error {
	if amount > 0 {
		deposits <- amount
		return nil
	}
	return fmt.Errorf("invalid amount")
}

func BalanceSafe() int {
	return <-balances
}

func Withdraw(amount int) bool {
	newWithdrawal := withdrawal{amount: amount, succeed: make(chan bool)}
	withdrawals <- newWithdrawal

	return <-newWithdrawal.succeed
}

// teller() is the monitor for variable "balance"
func teller() {
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case withdrawal := <-withdrawals:
			if withdrawal.amount < balance {
				balance -= withdrawal.amount
				withdrawal.succeed <- true
				continue
			}
			withdrawal.succeed <- false
		}

	}
}

type Cake struct{ state string }

func baker(cooked chan<- *Cake) {
	for {
		newCake := new(Cake)
		newCake.state = "cooked"
		// baker never touches this cake again
		cooked <- newCake
	}
}

func icer(iced chan<- *Cake, cooked <-chan *Cake) {
	for cake := range cooked {
		cake.state = "iced"
		// icer never touches this cake again
		iced <- cake
	}
}
