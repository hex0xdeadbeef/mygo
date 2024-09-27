package chapter5

import (
	"fmt"
	"net/http"

	// "net/http"
	"strings"

	"golang.org/x/net/html"
)

/* Basics */
func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }

func FunctionAssigning() {

	foo := square
	fmt.Println(foo(3))

	bar := negative
	fmt.Println(bar(10))

	zoo := product
	fmt.Println(zoo(10, 200))

	// foo = zoo /* cannot use zoo (variable of type func(m int, n int) int) as func(n int) int value in assignment */

	fmt.Printf("%T %T %T\n", foo, bar, zoo)
}

func NilFunction() {
	var foo func(int) int
	// foo(3) /* "runtime error: invalid memory address or nil pointer dereference" */
	fmt.Println(foo) /* nil */
}

func NilFunctionComparing() {
	var bar func(int, int) int

	if bar == nil {
		fmt.Println("bar is nil")
	}
}

func BehaviorChanging(str string) {
	fmt.Println(strings.Map(incrementRune, str))
}
func incrementRune(r rune) rune {
	return r + 1
}

/* Recursive HTML document traversals */
func TraversalOfHTML(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("requesting %s; %s", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", response.StatusCode)
	}

	node, err := html.Parse(response.Body)
	if err != nil {
		return err
	}
	id := "my-input"

	result := applyToEachNode(node, id, preID, nil)
	fmt.Println(result)
	return nil
}

var (
	depth int
)

func applyToEachNode(node *html.Node, id string, pre, post func(*html.Node, string) bool) *html.Node {
	if preID(node, id) {
		return node
	}

	if node.FirstChild != nil {
		result := applyToEachNode(node.FirstChild, id, preID, nil)
		if result != nil {
			return result
		}
	}

	if node.NextSibling != nil {
		result := applyToEachNode(node.NextSibling, id, preID, nil)
		if result != nil {
			return result
		}
	}

	return nil
}

func startElement(node *html.Node) {
	if node.Type == html.ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", node.Data)
		depth++
	}
}
func endElement(node *html.Node) {
	if node.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s<%s>\n", depth*2, "", node.Data)
	}
}

func headElement(node *html.Node) {
	if node.Type == html.ElementNode {
		depth++
		fmt.Printf("%*s%s\n", depth*2, "", node.Data)
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				fmt.Printf("%*s%s %s\n", depth*2, "", attr.Key, attr.Val)
			}
		}
	}
	if (node.Type == html.TextNode ||
		node.Type == html.CommentNode ||
		node.Type == html.DoctypeNode) && strings.Trim(node.Data, "\t\n\v\f\r \u0085\u00A0") != "" {
		fmt.Printf("%*s%s\n", depth*2, "", node.Data)
	}
}
func tailElement(node *html.Node) {
	if node.Type == html.ElementNode {
		fmt.Printf("%*s%s\n", depth*2, "", node.Data)
		depth--
	}
}

func preID(node *html.Node, id string) bool {
	if node.Type == html.ElementNode {
		if len(node.Attr) != 0 {
			for _, attr := range node.Attr {
				fmt.Println(attr)
				if attr.Key == "id" && attr.Val == id {
					return true
				}
			}
		}
	}
	return false
}

/* Multiple logic possibility */
func Placeholder(str string, handler func(string) string) string {
	portions := strings.Split(str, " ")
	for i, portion := range portions {
		if strings.HasPrefix(portion, "$") {

			portions[i] = handler(strings.TrimLeft(portion, "$"))
		}
	}

	return strings.Join(portions, " ")
}
func Replacements(mold string) string {
	switch mold {
	case "firstname":
		return "Dmitriy"
	case "lastname":
		return "Mamykin"
	case "age":
		return "21"
	default:
		return "empty"
	}
}
