package chapter8

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func SaveRequestBody(fileName string, url ...string) error {
	// Open file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating file with name \"%s\": %s", fileName, err)
	}
	defer file.Close()

	// Get the fastest response
	response, err := getFastestResponse(url)
	if err != nil {
		return fmt.Errorf("getting response: %s", err)
	}
	defer response.Body.Close()

	// Error handling
	_, err = io.Copy(file, response.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("copying data: %s", err)
	}
	if err := file.Close(); err != nil {
		if errors.Is(err, os.ErrClosed) {
			return fmt.Errorf("closing file: %s", err)
		} else {
			return fmt.Errorf("copying data: %s", err)
		}
	}

	return nil
}

type requestData struct {
	url     string
	context context.Context
}

func getFastestResponse(url []string) (*http.Response, error) {
	const (
		timeoutInSeconds = 3
	)

	var (
		wg          sync.WaitGroup
		responses   = make(chan *http.Response, len(url))
		ctx, cancel = context.WithCancel(context.Background())

		ticker = time.NewTicker(1 * time.Second)
	)
	defer func() {
		cancel()
		ticker.Stop()
	}()

	for _, url := range url {
		reqData := &requestData{url: url, context: ctx}
		wg.Add(1)
		go func() {
			defer wg.Done()
			request(reqData, responses)
		}()
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	for i := timeoutInSeconds; i > 0; i-- {
		select {
		case <-ticker.C:
		case response := <-responses:
			cancel()
			return response, nil
		}
	}
	ticker.Stop()

	return nil, fmt.Errorf("timeout")
}

func request(reqData *requestData, responses chan<- *http.Response) {
	request, err := http.NewRequestWithContext(reqData.context, "GET", reqData.url, nil)
	if err != nil {
		log.Printf("creating request to %s: %s", reqData.url, err)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("requesting \"%s\": resuest canceled", reqData.url)
			return
		} else {
			log.Printf("requesting \"%s\": %s", reqData.url, err)
			return
		}
	}

	// Goroutine's block won't happen because it's buffered channel.

	responses <- response

}
