package chapter8

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func SaveRequestBody(fileName string, number int, url ...string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating file with name \"%s\": %s", fileName, err)
	}

	response := getFastestResponse(number, url)
	defer response.Body.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		err = fmt.Errorf("reading response body of %s: %s", url, err)
	}
	file.Close()

	if err != nil {
		return err
	}

	return nil
}

func getFastestResponse(number int, url []string) *http.Response {
	responses := make(chan *http.Response, len(url))

	for _, url := range url {
		go request(url, responses)
	}

	return <-responses
}

func request(url string, responses chan<- *http.Response) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	responses <- response
}
