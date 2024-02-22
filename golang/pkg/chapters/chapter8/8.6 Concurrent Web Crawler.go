package chapter8

import (
	"context"
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
)

const (
	greedGoroutinesNumber = 32
)

var (
	tokens       = make(chan struct{}, greedGoroutinesNumber)
	traverseDone = make(chan struct{})
)

// AlternativelyModifiedCrawler() gets input from the os.Args() and gathers all the links within depth argument.
func AlternativelyModifiedCrawler(depth, goroutinesNumber int) map[string]bool {
	// Pair is used in depth limiting.
	type Pair struct {
		depth int
		links []string
	}

	var (
		maxDepth              = depth
		eagerGoroutinesNumber = goroutinesNumber

		wg               sync.WaitGroup
		workPairs        = make(chan *Pair)
		unseenLinksPairs = make(chan *Pair)
		ctx, cancel      = context.WithCancel(context.Background())

		seenLinks = make(map[string]bool)

		passedWorkersNumber int
	)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startPair := &Pair{depth: 1, links: os.Args[1:]}
		workPairs <- startPair
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		os.Stdin.Read(make([]byte, 1))

		cancel()
		close(traverseDone)
	}()

	for i := 0; i < eagerGoroutinesNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for unseenLinkPair := range unseenLinksPairs {
				if isTraverseCancelled() {
					return
				}

				// After passing to the "unseenLinksPairs" channel a new pair the pair has only on link at 0 index
				foundLinks, err := Crawler(ctx, unseenLinkPair.links[0])
				if err != nil {
					if urlErr, ok := err.(*url.Error); ok {
						if urlErr.Err == context.Canceled {
							log.Printf("\"%s\": request cancelled", unseenLinkPair.links[0])
							return
						}
					} else {
						log.Printf("crawling \"%s\": %s;", unseenLinkPair.links[0], err)
						continue
					}
				}

				if len(foundLinks) != 0 {
					newWorkPair := &Pair{depth: unseenLinkPair.depth + 1, links: foundLinks}
					wg.Add(1)
					go func() {
						defer wg.Done()
						workPairs <- newWorkPair
					}()
				}

			}
		}()
	}

	go func() {
		wg.Wait()
		close(workPairs)
		close(unseenLinksPairs)
	}()

loop:
	for {
		select {
		case <-traverseDone:
			return seenLinks

		case curWorkPair, ok := <-workPairs:
			if !ok {
				break loop
			}

			if curWorkPair.depth > maxDepth {
				continue
			}

			for _, link := range curWorkPair.links {
				if !seenLinks[link] {
					seenLinks[link] = true

					newUnseenLinksPair := &Pair{depth: curWorkPair.depth + 1, links: []string{link}}

					passedWorkersNumber++
					unseenLinksPairs <- newUnseenLinksPair
				}
			}
		}
	}

	return seenLinks
}

func isTraverseCancelled() bool {
	select {
	case <-traverseDone:
		return true
	default:
		return false
	}
}

func ResourceCopy(parsedLinks map[string]bool) error {
	const (
		directoryPath = "/Users/dmitriymamykin/Desktop/"
	)

	var (
		fileDescriptorsSemaphore = make(chan struct{}, 32)
	)

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
				log.Printf("copying page (%s) content: %s\n", url, err)
			}

			if err = websiteCopy(url, file); err != nil {
				log.Printf("copying page (%s) content: %s\n", url, err)
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

	links, err := ExtractLinksModified(nil, url)

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
