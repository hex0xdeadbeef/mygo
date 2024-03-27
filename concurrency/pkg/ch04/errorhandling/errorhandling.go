package errorhandling

import (
	"fmt"
	"net/http"
)

/*
In general, our concurrent processes should send theit errors to another part of our program that has complete information about the state of our program, and can make a
more informed decision about what to do.
*/

func NoErrorHandling() {
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)

		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}

				select {
				case <-done:
					return
				case responses <- resp:
				}

			}
		}()

		return responses
	}

	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

func ErrorHandling() {
	type RequestResult struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan struct{}, urls ...string) <-chan RequestResult {
		results := make(chan RequestResult)

		go func() {
			defer close(results)
			var reqRes RequestResult

			for _, url := range urls {
				reqRes = RequestResult{}
				reqRes.Response, reqRes.Error = http.Get(url)
				select {
				case <-done:
					return
				case results <- reqRes:
				}
			}
		}()

		return results
	}

	done := make(chan struct{})
	defer close(done)

	var errCounter int
	urls := []string{"https://www.google.com", "https://badhost", "a", "b", "c"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("requesting: %v\n", result.Error)
			errCounter++
			if errCounter > 2 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
