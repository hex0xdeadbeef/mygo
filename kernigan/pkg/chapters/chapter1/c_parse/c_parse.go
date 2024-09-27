package c_parse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func Fetch_Concurrently(filename string, urls ...string) {
	start := time.Now() // Start timer
	urlChannel := make(chan string)

	file, err := os.Create(filename)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
		os.Exit(1)
	}
	defer file.Close()

	for _, url := range urls {
		if !isHTTPPrefixed(url) {
			addHTTPPrefix(&url)
		}
		go fetch_Uploading_Channel(url, urlChannel)
	}

	for range urls {
		file.WriteString(fmt.Sprintf("%s\n", <-urlChannel))
	}
	file.WriteString(fmt.Sprintf("%.2fs elapsed\n", time.Since(start).Seconds()))
}

func fetch_Uploading_Channel(url string, channel chan<- string) {
	start := time.Now()

	response, err := http.Get(url)
	if err != nil {
		channel <- fmt.Sprintf(err.Error())
		return
	}

	numberOfWrittenBytes, err := io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		channel <- fmt.Sprintf("While reading %s: %s", url, err.Error())
	}

	secs := time.Since(start).Seconds()
	channel <- fmt.Sprintf("%.2fs %d %s", secs, numberOfWrittenBytes, url)
}

func Fetch(urls ...string) {
	for _, url := range urls {
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		numberOfWrittenBytes, err := io.Copy(os.Stdout, response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bytesArray := make([]byte, 0)
		bytesArray = append(bytesArray, byte(numberOfWrittenBytes))
		bytesArray = append(bytesArray, []byte(" "+response.Status+"\n")...)

		_, err = os.Stdout.Write(bytesArray)
		if err != nil {
			fmt.Println()
			os.Exit(1)
		}
	}

}

func isHTTPPrefixed(url string) bool {
	return strings.HasPrefix(url, "http://")
}

func addHTTPPrefix(url *string) {
	*url = "http://" + *url
}

func parsingWebsitesJson(fileName string) []string {
	type SiteData struct {
		Position int     `json:"position"`
		Domain   string  `json:"domain"`
		Count    int32   `json:"count"`
		Etv      float64 `json:"etv"`
	}

	file, err := os.Open("ranked_domains.json")
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
	}

	var jsonData []SiteData

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
	}

	var urls []string
	for i := 0; i < 100; i++ {
		urls = append(urls, jsonData[i].Domain)
	}

	return urls
}
