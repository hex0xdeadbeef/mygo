package main

import (
	"fmt"
	"sync"
)

type (
	User struct {
		ID int
	}

	Semaphore struct {
		C chan struct{}
	}

	resultWithErr struct {
		User User
		Err  error
	}
)

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.C
}

func (usr *User) Deactivate() error {
	return nil
}

func DeactivateUsersA(usrs []User, gCount int) ([]User, error) {
	var (
		wg sync.WaitGroup

		sem = Semaphore{
			C: make(chan struct{}, gCount),
		}

		outCh  = make(chan resultWithErr, len(usrs))
		sgnlCh = make(chan struct{})

		output = make([]User, 0, len(usrs))
	)

	for _, usr := range usrs {

		wg.Add(1)
		go func(usr User) {
			sem.Acquire()

			defer func() {
				sem.Release()
				wg.Done()
			}()

			err := usr.Deactivate()

			select {
			case outCh <- resultWithErr{User: usr, Err: err}:

			case <-sgnlCh:
				return

			}
		}(usr)

	}
	wg.Wait()

	for res := range outCh {
		if res.Err != nil {
			close(sgnlCh)
			return nil, fmt.Errorf("an error occured: %w", res.Err)
		}

		output = append(output, res.User)
	}

	return output, nil
}

func deactivateUser(producerCh <-chan User, consumerCh chan<- resultWithErr) {
	for usr := range producerCh {
		err := usr.Deactivate()

		consumerCh <- resultWithErr{User: usr, Err: err}
	}
}

func DeactivateUsersB(usrs []User, wgCnt int) ([]User, error) {
	var (
		producerCh = make(chan User)
		consumerCh = make(chan resultWithErr)

		wg = &sync.WaitGroup{}

		output = make([]User, 0, len(usrs))
	)

	// Producer goroutine
	go func() {
		defer close(producerCh)

		for _, usr := range usrs {
			producerCh <- usr
		}
	}()

	go func() {
		defer close(consumerCh)

		for i := 0; i < wgCnt; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				go deactivateUser(producerCh, consumerCh)
			}()

		}

		wg.Wait()
	}()

	for res := range consumerCh {
		if res.Err != nil {
			return nil, fmt.Errorf("an error occured: %w", res.Err)
		}

		output = append(output, res.User)
	}

	return output, nil

}
