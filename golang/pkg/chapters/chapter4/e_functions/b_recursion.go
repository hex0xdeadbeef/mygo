package e_functions

import (
	"fmt"
	sl "golang/pkg/chapters/chapter4/b_slices"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
	h "golang.org/x/net/html"
)

const (
	SpaceCutSet = "\t\n\v\f\r \u0085\u00A0"
)

/* Recursion */
// Print only links
func FindLinks1(url string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("resp status code: %v", response.StatusCode)
	}

	parsedBody, err := h.Parse(response.Body)
	if err != nil {
		return err
	}

	result := visit(nil, parsedBody)
	fmt.Println(result)
	return nil
}

func visit(imageLinkToCount map[string]int, node *html.Node) map[string]int {
	if imageLinkToCount == nil {
		imageLinkToCount = make(map[string]int)
	}
	if node.Type == html.ElementNode && isValidData(node) {
		imageLinkToCount[extractLink(node)]++
	}

	if node.FirstChild != nil {
		visit(imageLinkToCount, node.FirstChild)
	}

	if node.NextSibling != nil {
		visit(imageLinkToCount, node.NextSibling)
	}
	return imageLinkToCount
}
func isValidData(node *html.Node) bool {
	switch node.Data {
	// case "a", "img", "script", "style":
	case "img":
		return true
	}
	return false
}
func extractLink(node *html.Node) string {
	for _, attr := range node.Attr {
		if isValidKey(attr.Key) {
			return attr.Val
		}
	}
	return ""
}
func isValidKey(attr string) bool {
	switch attr {
	// case "src", "href":
	case "src":
		return true
	}
	return false
}

// Print only parent's/child's tags structure
func OutlinePrint() {
	file, err := os.OpenFile("body.html", os.O_RDONLY, os.ModeIrregular)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "OutlineUse: %v\n", err)
		os.Exit(1)
	}

	outline(nil, doc)
}
func outline(stack []string, node *html.Node) {
	if node.Type == html.ElementNode {
		stack = append(stack, node.Data) // Push current tag
		fmt.Println(stack)
	}

	if node.FirstChild != nil {
		outline(stack, node.FirstChild)
	}

	if node.NextSibling != nil {
		outline(stack, node.NextSibling)
	}
	fmt.Println(stack)
}

// Print a map of tag/count pairs
func PrintCounts(url string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("resp status code: %v", response.StatusCode)
	}

	parsedBody, err := h.Parse(response.Body)
	if err != nil {
		return err
	}

	counts := tagToCountMapping(nil, parsedBody)
	fmt.Println(counts)
	return nil
}
func tagToCountMapping(counts map[string]int, node *html.Node) map[string]int {
	if counts == nil {
		counts = make(map[string]int)
	}

	if node != nil {
		if node.Type == html.ElementNode {
			counts[node.Data]++
		}
	}
	if node.FirstChild != nil {
		tagToCountMapping(counts, node.FirstChild)
	}

	if node.NextSibling != nil {
		tagToCountMapping(counts, node.NextSibling)
	}

	return counts
}

// Print only text blocks
func OnlyTextPrint(url string) (int, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return -1, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("resp status code: %v", response.StatusCode)
	}

	parsedBody, err := h.Parse(response.Body)
	if err != nil {
		return -1, err
	}

	total := 0
	counts := textPrinting(nil, parsedBody)
	for key, value := range counts {
		total += value * len(strings.Split(string(sl.SquashSpacesSimple([]byte(key))), " "))
	}

	return total, nil
}
func textPrinting(stringCount map[string]int, node *html.Node) map[string]int {
	if stringCount == nil {
		stringCount = make(map[string]int)
	}

	if node.Parent != nil && node.Type == h.TextNode && !isScriptOrStyle(node.Parent) {
		stringCount[node.Data]++
	}

	if node.NextSibling != nil {
		textPrinting(stringCount, node.NextSibling)
	}

	if node.FirstChild != nil {
		textPrinting(stringCount, node.FirstChild)
	}
	return stringCount
}
func isScriptOrStyle(node *html.Node) bool {
	switch node.Data {
	case "script", "style":
		return true
	}
	return false
}
