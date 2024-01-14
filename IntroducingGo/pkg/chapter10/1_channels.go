package chapter10

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func pinger(channel chan<- string) {
	for {
		time.Sleep(time.Second * 1)
		channel <- "ping"
	}
}

func ponger(channel chan<- string) {
	for {
		time.Sleep(time.Second * 1)
		channel <- "pong"
	}
}

func printer(channel <-chan string) {
	for {
		message := <-channel
		fmt.Println("Hi, ", message)
	}
}

func Channel_Using() {
	var channel chan string = make(chan string)
	go pinger(channel)
	go printer(channel)
	go ponger(channel)

	var input string
	fmt.Scanln(&input)
}

// Channel that could get its own respective resourse first, writes data into it
func Select_Using() {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go func() {
		for {
			channel1 <- "From 1"
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			channel2 <- "From 2"
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case message1 := <-channel1:
				fmt.Println(message1)
			case message2 := <-channel2:
				fmt.Println(message2)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}

// Channel that could get its own respective resourse first, writes data into it
// But if during the time in argument After() (it's included in 3rd case) no one channel upload a datum,
// all channels sleep
func Timeout_Select_Using() {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go func() {
		for {
			channel1 <- "Data1"
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			channel2 <- "Data2"
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case message1 := <-channel1:
				fmt.Println("Message 1:", message1)
			case message2 := <-channel2:
				fmt.Println("Message 2:", message2)
			case <-time.After(time.Second * 1):
				fmt.Println("timeout")
				// default:
				// 	fmt.Println("nothing ready")
				// 	time.Sleep(time.Second * 10)
			}

		}
	}()

	var input string
	fmt.Scanln(&input)

}

type HomePage struct {
	URL  string
	Size int
}

func Buffered_Channel_Using() {
	urls := []string{
		"http://www.apple.com",
		"http://www.amazon.com",
		"http://www.google.com",
		"http://www.microsoft.com",
	}

	resultChannel := make(chan HomePage)
	for _, url := range urls {
		go func(url string) {
			// Get an answer from resource
			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			// Read info as bytes
			bs, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			// Put data into buffered channel
			resultChannel <- HomePage{
				URL:  url,
				Size: len(bs),
			}
		}(url)
	}

	var biggest HomePage

	for range urls {
		result := <-resultChannel
		if result.Size > biggest.Size {
			biggest = result
		}
	}

	fmt.Println("The biggest home page:", biggest.URL)

}
