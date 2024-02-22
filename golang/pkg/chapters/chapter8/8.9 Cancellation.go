package chapter8

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	done = make(chan struct{})
)

func cancelled() bool {
	select {
	// When channel "done" is closed this case is executed, because it gets zero values and respectively returns "true"
	case <-done:
		return true
	// When no byte is gotten, this function sends false
	default:
		return false
	}
}

func walkDirCancellation(curDir string, dirData *dirInfo, chanPair *channelPair) {
	for _, entry := range dirEntriesCancellation(curDir) {
		if entry.IsDir() {
			subDir := filepath.Join(curDir, entry.Name())
			chanPair.wg.Add(1)
			go func() {
				defer chanPair.wg.Done()
				/*
					When we get "true" the goroutine gets "true" and stops its execution respectively decrementing "wg" counter
				*/
				if cancelled() {
					return
				}
				walkDirCancellation(subDir, dirData, chanPair)
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

func dirEntriesCancellation(dir string) []fs.DirEntry {
	defer func() {
		<-semaphore
	}()
	select {
	case <-done:
		return nil
	case semaphore <- struct{}{}:
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("du: %s\n", err)
		return nil
	}

	return entries
}

func DiskUsageCancellation() {
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
		/*
			When we close the channel, the zero values are sent on the channel, so cancelled starts getting values and returns "true" values after it's called somewhere.
		*/
		close(done)
	}()

	go func() {
		wg.Wait()
		close(rootInfos)
	}()

loop:
	for {
		select {
		// Always get zero values after "done" has been drained
		case <-done:
			// When all the values that were sent on the channel "rootInfos" reached the main goroutine we'll return from the function
			for range rootInfos {
			}
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
