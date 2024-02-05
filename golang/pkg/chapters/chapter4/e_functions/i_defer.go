package e_functions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/net/html"
	. "golang.org/x/net/html"
)

/* Basics */
func Title(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("requesting to %s; %s", url, err)
	}
	/*
		Defer the call of closing the resourse
	*/
	defer response.Body.Close()

	// Get and check content type of the response
	contentType := response.Header.Get("Content-Type")
	if contentType != "text/html" && !strings.HasPrefix(contentType, "text/html") {
		return fmt.Errorf("%s has type %s, not text/html", url, contentType)
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	preProcess := func(node *Node) {
		if node.Type == ElementNode &&
			node.Data == "title" &&
			node.FirstChild != nil {
			fmt.Println(node.FirstChild.Data)
		}
	}

	forEachNode(doc, preProcess, nil)

	return nil
}

func BigSlowOperation() {
	//
	defer trace("BigSlowOperation")()
	/*
		... do some work
	*/
	time.Sleep(time.Second * 3)
}

func trace(message string) func() {
	start := time.Now()
	log.Printf("enter: %s", message)
	return func() { log.Printf("exit %s (%s)", message, time.Since(start)) }
}

func ObserveFunctionResult() {
	_ = doubleDeferred(2)
}

// After return the info about parameters and result will be printed
func doubleDeferred(x int) (result int) {
	defer func() { fmt.Printf("The argument of double(%d) = %d\n", x, result) }()
	result = 2 * x
	return
}

/*
After return statement but before an assigment the result variable will be incremented by the value of x
variable
*/
func Triple(x int) (result int) {
	defer func() { result += x }()
	result = double(x)
	return
}

func double(x int) int {
	return x + x
}

func DeferIntoForLoop() {
	for i := 0; i < 10; i++ {
		defer func() { fmt.Println(i) }()
	}
}

func DefeatFileDescriptorsExhaustion(filenames []string) error {
	for _, filename := range filenames {
		if err := workWithFile(filename); err != nil {
			return fmt.Errorf("working with: %s; %s", filename, err)
		}
	}

	return nil
}

func workWithFile(filename string) error {
	fileDescriptor, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening: %s; %s", filename, err)
	}
	defer fileDescriptor.Close()

	/*
		... some work
	*/
	return nil
}

func FetchImproved(url string) (filename string, n int64, err error) {
	response, err := http.Get(url)
	if err != nil {
		return "", -1, fmt.Errorf("while requesting: %s; %s", url, err)
	}
	/*
		Defer response body closing
	*/
	defer response.Body.Close()

	localPath := path.Base(response.Request.URL.Path)
	if localPath == "." || localPath == "/" {
		localPath = "index.html"
	}

	// Open a file
	fileDescriptor, err := os.Create(localPath)
	if err != nil {
		return "", -1, fmt.Errorf("creating %s; %s", localPath, err)
	}

	// Copy the info from HTML body into the file, any copying errors are postponed until the file closing.
	numberOfWrittenBytes, err := io.Copy(fileDescriptor, response.Body)
	defer func() {
		if err == nil {
			err = fmt.Errorf("closing %s; %s", localPath, fileDescriptor.Close())
			return
		}
		err = fmt.Errorf("copying %s HTML body to %s; %s; closing %s; %s", url, localPath, err, localPath, fileDescriptor.Close())
	}()

	return localPath, numberOfWrittenBytes, err
}
