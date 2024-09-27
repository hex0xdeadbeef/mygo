package chapter8

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileDescriptors = 32
)

var (
	semaphore = make(chan struct{}, maxFileDescriptors)
	mu        sync.Mutex
)

type dirInfo struct {
	root string

	filesNumber int64
	bytesNumber int64
}

type channelPair struct {
	c  chan<- *dirInfo
	wg *sync.WaitGroup
}

func walkDir(curDir string, dirData *dirInfo, chanPair *channelPair) {
	for _, entry := range dirEntries(curDir) {

		if entry.IsDir() {
			subDir := filepath.Join(curDir, entry.Name())
			chanPair.wg.Add(1)
			go func() {
				defer chanPair.wg.Done()
				walkDir(subDir, dirData, chanPair)
			}()

		} else {
			file, err := entry.Info()
			if err != nil {
				log.Printf("getting info about \"%s\": %s", file.Name(), err.Error())
				continue
			}

			mu.Lock()
			dirData.filesNumber++
			dirData.bytesNumber += file.Size()
			mu.Unlock()

			chanPair.c <- dirData
		}
	}

}

func dirEntries(dir string) []fs.DirEntry {
	defer func() {
		<-semaphore
	}()

	semaphore <- struct{}{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("du: %s\n", err)
		return nil
	}

	return entries
}

func DiskUsage() {
	const (
		tickersDurationInMillis = 500
	)

	var (
		verbose = flag.Bool("v", false, "show verbose progress messages")
		roots   []string

		wg        sync.WaitGroup
		rootInfos = make(chan *dirInfo)
		ticker    = time.NewTicker(1 * time.Millisecond)

		rootInfoCopy dirInfo
	)

	// Sets ticker's channel to the nil so that the select operator cannot select it
	ticker.C = nil

	flag.Parse()
	if *verbose {
		defer ticker.Stop()
		ticker = time.NewTicker(tickersDurationInMillis * time.Millisecond)
	}

	roots, err := getRoots()
	if err != nil {
		log.Print(err)

		roots = []string{"."}
		log.Print("directory set to the working directory")
	}

	chanPair := &channelPair{c: rootInfos, wg: &wg}
	for _, root := range roots {
		wg.Add(1)
		go func(initDir string) {
			defer wg.Done()
			rootData := &dirInfo{root: initDir}
			walkDir(initDir, rootData, chanPair)
		}(root)
	}

	go func() {
		os.Stdin.Read(make([]byte, 1))

		close(done)
	}()

	go func() {
		wg.Wait()
		close(rootInfos)
	}()

loop:
	for {
		select {
		case rootInfo, ok := <-rootInfos:
			if !ok {
				break loop
			}
			rootInfoCopy = *rootInfo
		case <-ticker.C:
			printDiskUsage(rootInfoCopy)
		}
	}
	printDiskUsage(rootInfoCopy)
}

func getRoots() ([]string, error) {
	if len(os.Args[1:]) == 0 {
		return nil, fmt.Errorf("getting prompt's args: no arguments provided")
	}

	return os.Args[1:], nil

}

func printDiskUsage(rootData dirInfo) {
	log.Printf("root \"%s\" has %d files %.3f GB\n",
		filepath.Base(rootData.root),
		rootData.filesNumber,
		float64(rootData.bytesNumber)/1e9)
}
