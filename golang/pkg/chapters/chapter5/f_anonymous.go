package e_functions

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Empty struct{}

/* Place-in anonymous function as an argument */
func Example() {
	fmt.Println(strings.Map(
		// This anonymous function will be passed into the function "strings.Map()"
		func(r rune) rune { return r + 1 },
		"HAL-9000"))
}

/*  */
func ReturnedFunctionUsing() {
	// The fucntions below have the different captured states.
	a := returnFunction()
	b := returnFunction()

	fmt.Println(a())
	fmt.Println(a())
	fmt.Println(a())

	fmt.Println(b())
}
func returnFunction() func() int {
	var x int

	// The function below captures the state of returnFunction()
	return func() int {
		x++
		return x * x
	}
}

/* Inner declaration and recursive using of an anonymous function */
func TopoSortUsing() {
	for i, course := range topologicalSort(prerequisites) {
		fmt.Printf("%d. %s\n", i, course)
	}
}

var (
	prerequisites = map[string][]string{
		"algorithms":            {"data structures"},
		"calculus":              {"linear algebra"},
		"compilers":             {"data structures", "formal languages", "computer organization"},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming", "discrete math"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}
)

func topologicalSort(graph map[string][]string) []string {
	order := make([]string, 0, len(graph))

	type Empty struct{}
	seen := make(map[string]Empty)

	var depthFirstSearch func(current string)

	depthFirstSearch = func(source string) {
		seen[source] = Empty{}

		for _, destination := range graph[source] {
			_, hasVisited := seen[destination]
			if !hasVisited {
				depthFirstSearch(destination)
				continue
			}
			fmt.Printf("%s --> %s\n", source, destination)
		}
		order = append(order, source)
	}

	for key := range graph {
		_, isPresented := seen[key]
		if !isPresented {
			depthFirstSearch(key)
		}
	}

	fmt.Println()
	return order
}

func extractLinksModified(url string) ([]string, error) {
	// Make a request and get response
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while requesting: %s; %s", url, err)
	}
	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server response code: %d", response.StatusCode)
	}
	// Get a root of HTML body
	doc, err := html.Parse(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("while parsing: %s; %s", url, err)
	}

	var links []string
	// The anonymous logic to extract valid links
	pre := func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key != "href" {
					continue
				}

				link, err := response.Request.URL.Parse(attr.Val)
				// Ignore the errors related to invalid relative/absolute links
				if err != nil {
					continue
				}
				// Append valid link into inner slice
				links = append(links, link.String())
			}
		}
	}

	// Traverse throughout all document and add all links
	forEachNode(doc, pre, nil)

	// Return links
	return links, nil
}
func forEachNode(node *html.Node, pre, post func(node *html.Node)) {
	if node != nil {
		pre(node)
	}

	if node.FirstChild != nil {
		forEachNode(node.FirstChild, pre, post)
	}
	if node.NextSibling != nil {
		forEachNode(node.NextSibling, pre, post)
	}

	if post != nil {
		post(node)
	}
}

func LinksBFS(sender func(string) []string, workspace []string) map[string]Empty {
	visited := make(map[string]Empty)

	// Unless all leaves are traversed
	for len(workspace) > 0 {
		copiedLinks := workspace
		workspace = nil

		for _, link := range copiedLinks {

			_, isVisited := visited[link]

			if !isVisited {
				visited[link] = Empty{}
				workspace = append(workspace, Crawler(link)...)
			}

		}
	}

	return visited
}

// Filter and return children links
func Crawler(url string) (links []string) {
	fmt.Println(url)
	// Get links
	links, err := extractLinksModified(url)

	// Gather links with only one hostname
	hostname := getHostname(url)
	filtredLinks := linkFilter(hostname, links)
	if err != nil {
		log.Printf("while extracting links from %s; %s", url, err)
	}

	return filtredLinks
}
func linkFilter(hostname string, sublinks []string) []string {
	validLinks := sublinks[:0]
	count := 0

	for _, sublink := range sublinks {
		subHostname := getHostname(sublink)
		if hostname == subHostname {
			validLinks = append(validLinks, sublink)
			count++
		}
	}
	return validLinks[:count]
}
func getHostname(rawUrl string) string {
	urlObj, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("while parsing %s; %s", rawUrl, err)
	}

	return urlObj.Hostname()
}

// Modified HTML elements printer
func Outline(url string) error {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("while requesting %s; %s", url, err)
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("after requesting %s; response status code: %d", url, response.StatusCode)
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("parsing %s; %s", url, err)
	}

	depth := 0

	preTagPrint := func(node *html.Node) {
		if node.Type == html.ElementNode {
			fmt.Printf("%*s%s\n", depth*2, "", node.Data)
			depth++
		}
	}

	postTagPrint := func(node *html.Node) {
		if node.Type == html.ElementNode {
			depth--
			fmt.Printf("%*s%s\n", depth*2, "", node.Data)
		}
	}

	forEachNode(doc, preTagPrint, postTagPrint)
	return nil
}
