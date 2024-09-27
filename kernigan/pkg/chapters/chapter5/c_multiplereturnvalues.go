package chapter5

import (
	"log"
	"net/http"

	h "golang.org/x/net/html"
)

func FindLinks2(url string) (map[string]int, error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("resp status code: %v", response.StatusCode)
	}

	parsedBody, err := h.Parse(response.Body)
	if err != nil {
		return nil, err
	}
	links := visit(nil, parsedBody)
	return links, nil
}

func CountWordsAndImages(url string) (words, images int, err error) {
	words, images, err = countWordsAndImages(url)
	return
}

func countWordsAndImages(url string) (int, int, error) {
	// Get images slice
	images, err := FindLinks2(url)
	if err != nil {
		log.Println(err)
		return -1, -1, err
	}
	imgsCount := len(images) - 1

	wrdsCount, err := OnlyTextPrint(url)
	if err != nil {
		log.Println(err)
		return -1, -1, err
	}
	return imgsCount, wrdsCount, nil
}
