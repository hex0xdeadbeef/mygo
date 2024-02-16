package chapter8

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Sequential processing
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

// With no waiting any goroutines
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go ImageFile(f)
	}
}

// Wait for all the goroutines
func makeThumbnails3(filenames []string) {
	done := make(chan bool)
	for _, f := range filenames {
		go func(f string) {
			ImageFile(f)
			done <- true
		}(f)
	}

	for range filenames {
		<-done
	}
}

// Goroutine's leaks
func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: goroutine leak!
		}
	}
	return nil
}

// Goroutine's leaks
func MakeThumbnails5(filenames []string) (thumbfiles []string) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))

	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			it.thumbfile = fmt.Sprintf("%s", it.err)
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}
	return thumbfiles
}

func MakeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // The number of working goroutines

	/*
		Placing wg.Wait() here would block the main goroutine and no children goroutines would be started
	*/

	for f := range filenames {
		/*
			Increment the current amount of workers
		*/
		wg.Add(1)
		// worker
		go func(f string) {
			/*
				Decrement the current amount of workers after all goroutine's activities done
			*/
			defer wg.Done()

			thumb, err := ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}

			info, _ := os.Stat(thumb)
			sizes <- info.Size()
		}(f) // Catch the current value f in new own mempool
	}

	/*
		Placing wg.Wait() here would lead to the infinite executing because there's no closer
	*/

	// closer
	go func() {
		/*
			Wait all the goroutines to finish
		*/
		wg.Wait()
		/*
			Close the size channel
		*/
		close(sizes)
	}()

	var total int64
	/*
		Wait for a new value of a size and add it to the total
	*/
	for size := range sizes {
		total += size
	}

	return total
}
