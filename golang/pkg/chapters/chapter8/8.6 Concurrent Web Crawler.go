package chapter8

import (
	"fmt"
	. "golang/pkg/chapters/chapter5"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func newLogger() (log.Logger, *os.File) {
	file, err := os.OpenFile("abc.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return *log.New(file, "-> ", log.Ldate|log.Ltime), file
}

const greedGoroutinesNumber = 50

var tokens = make(chan struct{}, greedGoroutinesNumber)

func ModifiedParallelLinksBFS() {
	worklist := make(chan []string)
	// Number of pending sends to worklist
	var n int

	// This statement is placed in its own goroutines because the goroutine would blocked after the operation
	// and no reading would complete. After the operation main goroutine will be blocked and no operations would execute.
	n++
	go func() {
		worklist <- os.Args[1:]
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			// Check whether the link was checked
			if !seen[link] {
				seen[link] = true
				// Launch the crawling the environment of the link that wasn't checked
				// and finally add its environment's links into the worklist
				n++
				go func(link string) {
					worklist <- ModifiedCrawler(link, seen)
				}(link)
			}
		}
	}

}

// Filter and return children links
func ModifiedCrawler(url string, seen map[string]bool) (links []string) {
	fmt.Println(url)

	// Only after one of the current goroutines releases the token other ones can access the expression
	tokens <- struct{}{}

	links, err := ExtractLinksModified(url)

	// When a goroutine has finished ExtractLinksModified() it releases the token and allows others to get the resourse
	<-tokens

	// Gather links with only one hostname
	hostname := GetHostname(url)
	filtredLinks := LinkFilter(hostname, links)
	if err != nil {
		log.Printf("while extracting links from %s; %s", url, err)
	}

	return filtredLinks
}

type Pair struct {
	depth int
	links []string
}

func AlternativelyModifiedCrawler(depth, goroutinesNumber int) map[string]bool {
	var (
		// Initialized with arguments
		maxDepth              = depth
		greedGoroutinesNumber = goroutinesNumber

		// Initialized
		workPair        = make(chan *Pair)
		unseenLinksPair = make(chan *Pair)

		seen = make(map[string]bool)

		timeout = 3

		// Uninitialized
		passedWorkersNumber int
	)

	defer func() {
		close(workPair)
		close(unseenLinksPair)
	}()

	go func() {
		initialPair := &Pair{depth: 1, links: os.Args[1:]}
		workPair <- initialPair
	}()

	for i := 0; i < greedGoroutinesNumber; i++ {

		// There will be 20 goroutines listen to the new unseenLinksPair
		go func() {

			for pairWithUnseenLinks := range unseenLinksPair {

				foundLinks := Crawler(pairWithUnseenLinks.links[0])
				newWorkPair := &Pair{depth: pairWithUnseenLinks.depth + 1, links: foundLinks}

				go func() {
					workPair <- newWorkPair
				}()
			}
		}()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := timeout; i > 0; i-- {
		select {

		case <-ticker.C:

		case curWorkPair := <-workPair:
			if curWorkPair.depth > maxDepth {
				continue
			}

			for _, link := range curWorkPair.links {
				if !seen[link] {
					seen[link] = true

					newUnseenLinksPair := &Pair{depth: curWorkPair.depth + 1, links: []string{link}}

					passedWorkersNumber++
					unseenLinksPair <- newUnseenLinksPair
				}
			}
			i += timeout
		}
	}

	return seen
}

func ResourceCopy(parsedLinks map[string]bool) error {
	const (
		directoryPath = "/Users/dmitriymamykin/Desktop/"
	)

	var (
		logger, loggerFile       = newLogger()
		fileDescriptorsSemaphore = make(chan struct{}, 32)
	)
	defer loggerFile.Close()

	var err error
	if len(parsedLinks) == 0 {
		return fmt.Errorf("links ")
	}

	var website string
	for key := range parsedLinks {
		if website, err = websiteName(key); err != nil {
			return fmt.Errorf("getting website name: %s", err)
		}
		break
	}

	if err := os.Mkdir(directoryPath+website, 0777); err != nil {
		return fmt.Errorf("creating directory: %s", err)
	}

	if err = setWorkDir(directoryPath + website); err != nil {
		return fmt.Errorf("changing work directory: %s", err)
	}

	var errorChan = make(chan error)
	var wg sync.WaitGroup

	for key := range parsedLinks {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			fileName := fileName(url)

			// Sets the limit on the number of concurrent launched file descriptors tokens.
			fileDescriptorsSemaphore <- struct{}{}
			file, err := fileDescriptor(fileName)
			<-fileDescriptorsSemaphore

			if err != nil {
				logger.Printf("copying page (%s) content: %s\n", url, err)
			}

			if err = websiteCopy(url, file); err != nil {
				logger.Printf("copying page (%s) content: %s\n", url, err)
			}
		}(key)

	}

	go func() {
		wg.Wait()
		close(errorChan)
		close(fileDescriptorsSemaphore)
	}()

	for err := range errorChan {
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

func fileName(url string) string {
	return filepath.Base(url) + ".html"
}

func fileDescriptor(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("creating file with name: %s: %s", fileName, err)
	}

	return file, nil
}

func setWorkDir(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}

	return nil
}

func websiteName(initialURL string) (string, error) {

	parsedURL, err := url.Parse(initialURL)
	if err != nil {
		return "", fmt.Errorf("parsing url: %s: %s", initialURL, err)
	}

	return strings.Split(parsedURL.Hostname(), ".")[1], nil
}

func websiteCopy(url string, file *os.File) error {
	response, err := http.Get(url)
	if err != nil {
		fmt.Errorf("requesting \"%s\": %s", url, err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil && err != io.EOF {
		return fmt.Errorf("while copying resource: %s", err)
	}
	if file == nil {
		return fmt.Errorf("File doesn't exist")
	}
	return file.Close()
}
