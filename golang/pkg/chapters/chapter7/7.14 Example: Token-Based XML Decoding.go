package chapter7

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
HOMEWORK
*/
// 7.17
func SelectXML(fileName string) error {
	// Open file with XML data
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err, ok := err.(*os.PathError); ok {
		return err.Err
	} else if err != nil {
		return err
	}
	defer file.Close()

	// Define and initialize the decoder based on the file above
	decoder := xml.NewDecoder(file)
	// stack of element names
	var stack []xml.StartElement
	// Get required start elements
	requiredStartElems, err := constructStartElems()
	if err != nil {
		return fmt.Errorf("creating required start elems; %s", err)
	}

	// While we don't reach the EOF
	for {
		// Get next token
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("getting next token; %s\n", err)
		}

		switch token := token.(type) {
		case xml.StartElement:
			// Put the name onto the stack
			stack = append(stack, token)
		case xml.EndElement:
			// Reduce the stack by cutting up the only one element from the bottom
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if containsAll(stack, requiredStartElems) {
				fmt.Printf("%s-> %s\n", stackJoin(stack), token)
			}
		}
	}

	return nil
}

func constructStartElems() ([]xml.StartElement, error) {
	stringStartElemsNames, err := getStartElemsNames()
	if err != nil {
		return nil, fmt.Errorf("getting start elems names; %s", err)
	}
	startElems, err := startElemsInstances(stringStartElemsNames)
	if err != nil {
		return nil, fmt.Errorf("creating a slice of named start elems instances; %s", err)
	}

	stringParamsAndValues, err := getParamsAndValues(startElems)
	if err != nil {
		return nil, fmt.Errorf("getting params and values for start elems: %s; %s", getStartElementsNamesString(startElems), err)
	}

	err = fillStartElemsWithParamsAndValues(startElems, stringParamsAndValues)
	if err != nil {
		return nil, fmt.Errorf("filling start elems slice with params and its values; %s", err)
	}

	return startElems, nil
}

// containsAll reports whether x contains the element of y in order.
func containsAll(x, y []xml.StartElement) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if areStartElemsEqual(x[0], y[0]) {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

func areStartElemsEqual(x, y xml.StartElement) bool {
	return startElementString(x) == startElementString(y)
}

func stackJoin(stack []xml.StartElement) string {
	result := ""
	for _, se := range stack {
		result += se.Name.Local + " "
	}
	return result
}

func startElementString(x xml.StartElement) string {
	result := ""
	result += x.Name.Local
	for _, attr := range x.Attr {
		result += attr.Name.Local + attr.Value
	}

	return result
}

func getStartElemsNames() ([]string, error) {
	const maxAttemptsCount = 5

	var (
		userInput    []string
		attemptCount int
	)

	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Print("Write attributes names separated by commas -> ")

	// Get values for parameters from user
	for scanner.Scan() {
		userInput = strings.Split(
			strings.Trim(scanner.Text(), "\t\n\v\f\r \u0085\u00A0"),
			",")

		attemptCount++

		if attemptCount == maxAttemptsCount {
			return nil, errors.New("Attempts limit exceeded")
		}

		fmt.Printf("The attributes \"%v\" has been initially accepted.\n", userInput)
		break
	}

	return userInput, nil
}

func getStartElementsNamesString(startElements []xml.StartElement) string {
	result := "| "
	for _, startElement := range startElements {
		result += (string(startElement.Name.Local) + " ")
	}
	result += "|"

	return result
}

func startElemsInstances(startElemsStr []string) ([]xml.StartElement, error) {
	startElems := make([]xml.StartElement, 0, len(startElemsStr))
	for _, startElementName := range startElemsStr {
		if startElementName != "" {
			startElems = append(startElems, xml.StartElement{Name: xml.Name{Local: startElementName}, Attr: []xml.Attr{}})
			continue
		}
		return nil, fmt.Errorf("invalid start element name; name: \"%s\"", startElementName)
	}

	return startElems, nil
}

func getParamsAndValues(startElements []xml.StartElement) ([][]string, error) {
	const maxAttemptsCount = 5

	var (
		userInput    []string
		attemptCount int
	)

	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Printf("Write attributes and its values for %s start elems in the next way attr11,val11,...,attr1N,val1N|... -> ", getStartElementsNamesString(startElements))

	// Get values for parameters from user
	for scanner.Scan() {
		userInput = strings.Split(
			strings.Trim(scanner.Text(),
				"\t\n\v\f\r \u0085\u00A0"),
			"|")

		attemptCount++

		if attemptCount == maxAttemptsCount {
			return nil, errors.New("Attempts limit exceeded")
		}

		if len(userInput) != len(startElements) {
			fmt.Printf("invalid values; values: \"%v\"\n", userInput)
			fmt.Printf("Write attributes and its values for %s start elems in the next way attr1,val1|...attrN,valN -> ", getStartElementsNamesString(startElements))
			continue
		}

		fmt.Printf("The values \"%v\" has been initially accepted.\n", userInput)
		break
	}

	vectors := make([][]string, 0, len(userInput))
	for _, vector := range userInput {
		vectors = append(vectors, strings.Split(vector, ","))
	}

	return vectors, nil
}

func fillStartElemsWithParamsAndValues(startElements []xml.StartElement, paramsAndValues [][]string) error {
	for i, vector := range paramsAndValues {
		if len(vector)%2 != 0 {
			return fmt.Errorf("insufficient input; %v", vector)
		}
		for j := 0; j < len(vector); j += 2 {
			startElements[i].Attr = append(startElements[i].Attr, xml.Attr{Name: xml.Name{Local: vector[j]}, Value: vector[j+1]})
		}
	}
	return nil
}

func printStartElemsSlice(startElems []xml.StartElement) {
	for _, se := range startElems {
		fmt.Println("----------------------------")
		fmt.Printf("Attr: %s\n", se.Name.Local)
		for _, attr := range se.Attr {
			fmt.Printf("%s -> %s\n", attr.Name.Local, attr.Value)
		}
		fmt.Println("----------------------------")
	}
}

func FetchModified(fileName string, urls ...string) error {
	for _, url := range urls {
		response, err := http.Get(url)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err, ok := err.(*os.PathError); ok {
			return err.Err
		} else if err != nil {
			return err
		}
		defer file.Close()

		// The io.Copy() postpones any errors until the file closed, so we must explicitly close the file and catch any errors.
		nwb, err := io.Copy(file, response.Body)
		response.Body.Close()
		if err != nil {
			return err
		}

		bytesArray := make([]byte, 0)
		bytesArray = append(bytesArray, byte(nwb))
		bytesArray = append(bytesArray, []byte(" "+response.Status+"\n")...)

		_, err = file.Write(bytesArray)
		if err != nil {
			return err
		}
	}

	return nil
}

// 7.18
type Node interface{}

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func ConstructTree(fileName string) (Node, error) {

	// Open file with XML data
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err, ok := err.(*os.PathError); ok {
		return nil, err.Err
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	// The head of document tree
	var (
		root     Node
		curToken xml.Token
		stack    []*Element
	)
	const (
		spaceCutSet = "\t\n\v\f\r \u0085\u00A0"
	)

	tokenDecoder := xml.NewDecoder(file)

	rootElement := Element{}
	root = Node(&rootElement)

	stack = append(stack, &rootElement)

	for {
		// Get next token
		curToken, err = tokenDecoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch curToken := curToken.(type) {
		case xml.StartElement:
			stack = startElemProcess(&curToken, stack)
		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("invalid xml data; attempt to reduce the stack empty with the end element %s", curToken.Name.Local)
			}

			if curToken.Name.Local != stack[len(stack)-1].Type.Local {
				return nil, fmt.Errorf("invalid xml data; start and end element names don't match; start element name: \"%s\", end element name: \"%s\"", curToken.Name.Local, stack[0].Type.Local)
			}

			if stack = endElementProcess(stack); len(stack) == 0 {
				break
			}

		case xml.CharData:
			if strings.Trim(fmt.Sprintf("%s", curToken), spaceCutSet) == "" {
				continue
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("invalid xml data; attempt to add the char data to the last element of the empty stack; char data: \"%s\"", curToken)
			}
			strData := CharData(fmt.Sprintf("%s", curToken))
			charDataProcess(strData, stack)
			fmt.Println()
		}
	}

	return root, nil
}

func startElemProcess(startElem *xml.StartElement, stack []*Element) []*Element {
	// Construct a new start elements
	element := Element{Type: startElem.Name, Attr: startElem.Attr, Children: []Node{}}

	var stacklLength int = len(stack)

	// Add a child to the last element that is on the top of the stack
	stack[stacklLength-1].Children = append(stack[stacklLength-1].Children, Node(&element))
	stack = append(stack, &element)
	return stack

}

func endElementProcess(stack []*Element) []*Element {
	// Reduce the stack
	stack = stack[:len(stack)-1]
	return stack
}

func charDataProcess(charData CharData, stack []*Element) {
	stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, Node(charData))
}

func PrintTree(root Node) {
	var depth int

	var traverse func(Node)

	traverse = func(n Node) {
		depth++
		switch n := n.(type) {
		case *Element:
			fmt.Printf("%*s%s ", depth*2, "", n.Type.Local)
			for _, attr := range n.Attr {
				fmt.Printf("%s ", attr.Value)
			}
			fmt.Println()
			for _, child := range n.Children {
				traverse(child)
			}
		case CharData:
			fmt.Printf("%*s%s\n", depth*2, "", n)
		default:
			fmt.Printf("%T\n", n)
		}

		depth--
	}

	traverse(root)
}
