package e_functions

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

/* Basics */
func Sum(vals ...int) int {
	vals[0] = 100

	total := 0
	for _, value := range vals {
		total += value
	}

	return total
}

func TypesDifference() {
	var sliceParameterFunc func([]int)
	var variadicParameterFunc func(...int)

	fmt.Printf("slice par func: %T\n", sliceParameterFunc)
	fmt.Printf("variadic par func: %T\n", variadicParameterFunc)
}

func Min(vals ...int) (int, error) {
	if len(vals) == 0 {
		return -1, fmt.Errorf("no required arguments")
	}

	min := 0
	for _, value := range vals {
		if value <= min {
			min = value
		}
	}

	return -1, nil
}

func Max(vals ...int) (int, error) {
	if len(vals) == 0 {
		return -1, fmt.Errorf("no required arguments")
	}

	max := 0
	for _, value := range vals {
		if value >= max {
			max = value
		}
	}

	return -1, nil
}

func VariadicJoin(separator string, strs ...string) (string, error) {
	if len(strs) == 0 {
		return "", fmt.Errorf("no strings provided")
	}

	return strings.Join(strs, separator), nil
}

func ElementsByTagName(doc *html.Node, names ...string) ([]*html.Node, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no names provided")
	}

	var elementsCount int
	var elements []*html.Node

	isNameValid := func(node *html.Node) bool {
		elementName := node.Data
		for _, name := range names {
			if elementName == name {
				return true
			}
		}
		return false
	}

	pre := func(node *html.Node) {
		if node.Type == html.ElementNode && isNameValid(node) {
			elements = append(elements, node)
			elementsCount++
		}
	}

	forEachNode(doc, pre, nil)
	return elements, nil
}

func GetTopHTMLNode(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while requesting %s; %s", url, err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("after request %s; response status code: %d", url, response.StatusCode)
	}

	topNode, err := html.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("while parsing %s; %s", url, err)
	}

	return topNode, nil
}
