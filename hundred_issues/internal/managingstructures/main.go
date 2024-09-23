package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

/*
	MANAGING OPERATORS ERRORS


	IGNORING ELEMS COPYNG DURING for-range LOOP
1. If we're not interested in using value in for-range loop we can skip the second variable.

	VALUES COPYING IN for-range LOOP
1. During each iteration the for-range form of looping copies the value.
2. During assingment we pass the copy of the source val.
	1) If we assign the result of function execution returning an instance of a struct, Go creates a copy of the instance
	2) If we assign the result of execution of the function returning a pointer, Go creates the copy of the address of the memory area. (In 64-bit architecture an
	address has the length of 64-bit)
3. During an iteration over a data structure (slice, map) the for-loop range creates a copy of the current val and assigns it to the second variable. In the case
when we want to alternate the value of the fields, we can:
	1) Directly refer to the elem with an index, using the classic form of for-loop
	2) Use the index from for-range loop to change slice elems
	3) Turn the field into the pointer type
	4) Turn the slice into pointers to objects
	This method isn't effiecient one, because it slows down the branch-prediction + we must change the type of the slice and it's not always possible to do.

	IGNORING THE WAY OF ARGUMENTS EVALUATION IN for-range LOOPS
1. The structure we will loop over is evaluated before the looping.
2. The basic form of for-loop is different from for-range loop. We can restructure the underlying data and it will refer to the current version of it.

The logic of copying during loop over slice is appliable to other types: channels, arraysю.
3. Working wigh channels the evaluation of the parameter to be looped over is evaluated before the execution of loop-statement, so the swap the variable to a new
value won't affect the instance we're looping over.
4. Before the looping over an array, this array will be copied to an internal variable. To observe the changes we made during the loop we can use:
	1) Direct reference to the source array by indexes
	2) Passing a pointer to the array. In this case the underlying array elems won't be copied, it results in the higher performance.

	IGNORING THE EFFECT OF USING POINTERS IN for-range LOOPS
1) Using for-range loop and assigning to the internal strucure's element a pointer to the loop's variable results in storing in the target structure only the last element. The
problem is storing this variable in its memory area. At each iteration we copy a new elem to this variable and assign to the target data structure element a pointer to it ->
finally we will store the pointer to this variable updated up to the last element.

	LOOPING OVER THE MAP
1) Map doesn't store elems in the sorted order (It's not based on the binary-tree)
2) Map doesn't store the order in which elems were be added into it.
3) We cannot assume what order elems will have after inserting it into a map
4) Updating the map during iterations.There's no guarantee that the element will be inserted into a map during the iterations over the map. The specification says:
	If the writing into a map is happening during the iteration, it can be executed or skipped at all. The choice can vary for each created insertion and from one iteration to
	another one.
To beat this behavior we must create a copy of the source map and make insertions into the copy, not the initial map.

	IGNORING PROPERTIES OF BREAK OPERATOR
1. There are a bunch of error while working with switch/select statements in combination with break operator. Programmers break the case, not the cycle
2. break operator breaks the closest to it statement from these ones: switch | select | for
3. An operator switch | select | for with the break operator is an idiomatic thing in Go
4. The same logic is applied for continue operator:
	1) continue continues the closest to it switch | select | for operator.

	USING DEFER IN CYCLES
1. defer plans the actions inside it after returning of the enclosing function
2. In the case of a restrictive circumstances (for example: file descriptors can be exhausted) we can apply the following patterns:
	0) Arrange the work with closing the files after each iteration of the cycle
	1) Enclose all the actions within a cycle in a function forcing the defer to work after return statement of this function. In this case the defer statement is guaranteed to
	be executed after each iteration. It results in an overhead.
	2) To make the func readFile(...) a closing function. In fact this solution is the same to 1). It results in an overhead too.
3. We should remember that in a cycle defer statements just are piled to a stack and will be executed after the return statement of the function, so we need to beat this problem

*/

func main() {
	// SkippingValInForRange()

	// ValueCopyingDuringForRange()

	// ArguentsEvaluationInLoops()

	// ArgumentsEvaluationInChans()

	// ArgumentsEvaluationInArrays()

	// PointersInForRangeLoop()

	// MapAdding()

	// InsertingIntoMap()

	// UpdatingMapDuringIteration()

	// BreakSwitchOrSelectWithCycle()

	// DeferInCycles()

	// workWithProdA()

	workWithProdB()
}

func SkippingValInForRange() {
	var (
		s = []string{"a", "b", "c"}
	)

	// i - is the index, not the value
	for i := range s {
		fmt.Println(i)
	}
}

func ValueCopyingDuringForRange() {
	var (
		a = func() {
			type (
				account struct {
					balance float64
				}
			)

			const (
				salaryRise = 1000
			)

			accs := []account{
				{balance: 100},
				{balance: 200},
				{balance: 300},
			}

			// v gets a copy of the current acc held by accs slice
			for _, v := range accs {
				// The incrementation of the v variable results in the incrementing the copy.
				v.balance += salaryRise
			}

			for _, v := range accs {
				fmt.Println(v)
			}
		}

		b = func() {
			type (
				balance struct {
					amount float64
				}

				account struct {
					b *balance
				}
			)

			const (
				salaryRise = 1000
			)

			accs := []account{
				{b: &balance{100}},
				{b: &balance{200}},
				{b: &balance{300}},
			}

			// v gets a copy of the current acc held by accs slice
			for _, v := range accs {
				// Since the v now stores a pointer to the b field, a copy in v variable will have the copy of the pointer -> the incrementation wull be viewed in
				// the future
				v.b.amount += salaryRise
			}

			for _, v := range accs {
				fmt.Println(v.b.amount)
			}
		}

		c = func() {
			type (
				account struct {
					balance float64
				}
			)

			const (
				salaryRise = 1000
			)

			accs := []account{
				{balance: 100},
				{balance: 200},
				{balance: 300},
			}

			// this form of for-loop allows us to use the index in a pleasant way
			for i := range accs {
				// The incrementation will be applied and the result of it will be reflected in the source
				accs[i].balance += salaryRise
			}

			for _, v := range accs {
				fmt.Println(v)
			}
		}

		d = func() {
			type (
				account struct {
					balance float64
				}
			)

			const (
				salaryRise = 1000
			)

			accs := []account{
				{balance: 100},
				{balance: 200},
				{balance: 300},
			}

			// this form of for-loop allows us to use the index in a pleasant way
			for i := 0; i < len(accs); i++ {
				// The incrementation will be applied and the result of it will be reflected in the source
				accs[i].balance += salaryRise
			}

			for _, v := range accs {
				fmt.Println(v)
			}
		}

		e = func() {
			type (
				account struct {
					balance float64
				}
			)

			const (
				salaryRise = 1000
			)

			accs := []*account{
				{balance: 100},
				{balance: 200},
				{balance: 300},
			}

			// this form of for-loop allows us to use the index in a pleasant way
			for _, v := range accs {
				// The incrementation will be applied and the result of it will be reflected in the source
				v.balance += salaryRise
			}

			for _, v := range accs {
				fmt.Println(v)
			}
		}
	)

	fmt.Println()
	a()

	fmt.Println()
	b()

	fmt.Println()
	c()

	fmt.Println()
	d()

	fmt.Println()
	e()

}

func ArgumentsEvaluationInLoops() {
	var (
		a = func() {
			var (
				s = []int{0, 1, 2}
			)

			// The current copy of the s will be evaluated before the looping over.
			// The temporary variable in the loop will have len = 3, cap = 3.
			for range s {
				// The cap will be expended after first append
				fmt.Println(len(s), cap(s))

				s = append(s, rand.Intn(128))
			}

			fmt.Println()

			for _, v := range s {
				fmt.Print(v, " ")
			}
		}

		b = func() {
			var (
				s = []int{0, 1, 2}
			)

			// The basic form of for-loop evaluates the actual state of the condition, so this cycle won't be broken at all.
			for i := 0; i < len(s); i++ {
				s = append(s, rand.Intn(128))
			}

			for _, v := range s {
				fmt.Println(v)
			}
		}
	)

	a()

	b()
}

func ArgumentsEvaluationInChans() {

	var (
		ch1 = make(chan int, 3)
	)

	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	var (
		ch2 = make(chan int, 3)
	)

	go func() {
		for i := 10; i < 13; i++ {
			ch2 <- i
		}

		close(ch2)
	}()

	ch := ch1

	// The channel we iterate over is evaluated before the execution of loop statement
	for v := range ch {
		fmt.Println(v)

		// The loop will resume looping over ch1, because the evaluation of looping data is calculated before the loop execution
		ch = ch2 // There's the effect of swapping the channel to the ch2. If we close the ch channel, the ch2 will be closed instead of ch1
	}
}

func ArgumentsEvaluationInArrays() {
	var (
		a = func() {
			var (
				arr = [3]int{1, 2, 3}
			)

			// The array is copied before the execution of a for-loop statement, but we can restructure the internals of the initial elems with the index "i"
			for i, v := range arr {
				arr[i] = rand.Intn(128)

				fmt.Println(i, v)
			}

			// The internal of the source elems will be changed
			fmt.Println(arr)
		}

		b = func() {
			type (
				customInternalInt struct {
					val int
				}
			)

			var (
				arr = [3]*customInternalInt{{1}, {2}, {3}}
			)

			// The copied array will be the array
			for i, v := range arr {
				v.val = rand.Intn(128)
				fmt.Println(arr[i])
			}

		}

		c = func() {
			var (
				arr = [3]int{1, 2, 3}
			)

			// Before the execution of the for-range loop the pointer to the underlying arr will be copied
			for i, v := range &arr {
				// The elems of the source arr will be alternated instead of the copies
				if idx := i + 1; idx < len(&arr) {
					arr[idx] = rand.Intn(128)
				}

				// After the first step we will see the updated vals by the previous steps
				fmt.Println(i, v)
			}
		}
	)

	fmt.Println()
	a()

	fmt.Println()
	b()

	fmt.Println()
	c()
}

type (
	Customer struct {
		ID      string
		Balance float64
	}

	Store struct {
		m map[string]*Customer
	}
)

func New() *Store {
	return &Store{m: make(map[string]*Customer, 1<<7)}
}

func (s *Store) storeCustomers(customers []Customer) {
	var (
		a = func() {
			// the loop's variable customer has the same address in memory, so when we assign to the map the address of this variable we store
			// the same address at each iteration
			for _, customer := range customers {
				// Shading the customer variable to store the different pointers
				cur := customer
				// Store the shading variable
				s.m[customer.ID] = &cur
			}
		}

		b = func() {
			for i := range customers {
				// We're referencing to the source datum (customer) with a current index
				s.m[customers[i].ID] = &customers[i]
			}
		}
	)

	if rand.Intn(2) == 1 {
		a()
		return
	}
	b()

}

func PointersInForRangeLoop() {
	s := New()

	s.storeCustomers([]Customer{
		{ID: "1", Balance: 10},
		{ID: "2", Balance: -10},
		{ID: "3", Balance: 0},
	})

	for _, v := range s.m {

		fmt.Println(v)
	}
}

func MapAdding() {
	var (
		m = make(map[string]struct{}, 1<<6)
	)

	for _, v := range []string{"a", "y", "z", "c", "d", "e"} {
		m[v] = struct{}{}
	}

	for k := range m {
		fmt.Println(k)
	}

	/*
		d
		c
		e
		y
		a
		z
	*/
}

func InsertingIntoMap() {
	var (
		m = map[int]bool{
			0: true,
		}
	)

	// The copying the map doesn't happen, so the number of additions isn't deterministic.
	/*
		The outputs are:
			map[0:true 10:true 20:true 30:true 40:true 50:true]
			map[0:true 10:true 20:true]
			map[0:true 10:true 20:true 30:true]


	*/
	for k, v := range m {

		if v {
			m[k+10] = true
		}
	}

	fmt.Println(m)
}

func UpdatingMapDuringIteration() {
	var (
		copyMap = func(m map[int]bool) map[int]bool {
			var (
				res = make(map[int]bool, len(m))
			)

			for k, v := range m {
				res[k] = v
			}

			return res
		}

		m = map[int]bool{
			0: true,
			1: false,
			2: true,
		}

		copiedM = copyMap(m)
	)

	for k, v := range m {
		if v {
			copiedM[k+10] = true
		}
	}

	fmt.Println(copiedM) // map[0:true 1:false 2:true 10:true 12:true]
}

func BreakSwitchOrSelectWithCycle() {
	var (
		a = func() {
		Out:
			for i := 0; i < 5; i++ {
				switch i {
				case 2:
					// This statement breaks the switch only, not the cycle
					break Out
				default:
					fmt.Println(i)
				}
			}
		}

		aFixed = func() {
		FirstLoop:
			for i := 0; i < 5; i++ {
				switch i {
				case 2:
					// We break the loop with marker "FirtsLoop" here
					break FirstLoop
				default:
					fmt.Println(i)
				}
			}

		}

		b = func() {
			var (
				producer = func(ctx context.Context) <-chan rune {
					producer := make(chan rune)

					go func() {
						defer close(producer)

						// The marker is responsible for being used while breaking
					firstLoop:
						for {
							select {
							case <-ctx.Done():
								// In this case we break the loop in general, not the select operator
								break firstLoop
							case producer <- rune(rand.Intn(128)):

							}
						}

					}()

					return producer
				}
			)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Microsecond*500))
			defer cancel()

			for r := range producer(ctx) {
				fmt.Printf("%c ", r)
			}

		}
	)

	a()
	/*
		0
		1
		3
		4
	*/

	fmt.Println()
	aFixed()
	/*
		0
		1
	*/

	fmt.Println()
	b()

}

func DeferInCycles() {
	if errs := readFilesA(nil); errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	if errs := readFilesB(nil); errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	if errs := readFilesС(nil); errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

}

func readFilesA(filePaths <-chan string) []error {
	if filePaths == nil {
		return []error{fmt.Errorf("filePaths channel is nil")}
	}

	var (
		errs = make([]error, 1<<8)
	)

	for fp := range filePaths {
		if err := readFileA(fp); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func readFileA(filePath string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", filePath, err)
	}
	defer func(filePath string) {
		fileClosingErr := f.Close()

		if fileClosingErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; closing file %q: %w", err, filePath, fileClosingErr)
			return
		}

		err = fmt.Errorf("closing file %s: %q", filePath, fileClosingErr)

	}(filePath)

	// some logic...

	return err
}

func readFilesB(filePaths <-chan string) []error {
	if filePaths == nil {
		return []error{fmt.Errorf("filePaths channel is nil")}
	}

	var (
		errs = make([]error, 1<<8)
	)

	for fp := range filePaths {
		err := func(filePath string) (err error) {
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening file %q: %w", filePath, err)
			}
			defer func(filePath string) {
				fileClosingErr := f.Close()

				if fileClosingErr == nil {
					return
				}

				if err != nil {
					err = fmt.Errorf("%w; closing file %q: %w", err, filePath, fileClosingErr)
					return
				}

				err = fmt.Errorf("closing file %s: %q", filePath, fileClosingErr)

			}(filePath)

			// some logic...

			return err
		}(fp)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

type (
	Semaphore struct {
		c chan struct{}
	}

	FileWorker struct {
		file     *os.File
		filepath string
	}

	ErrsPipe struct {
		errorsCh chan error
	}
)

func readFilesС(filePaths <-chan string) (errs []error) {
	const (
		fileDescriptorsLimit = 1 << 10
	)

	if filePaths == nil {
		return []error{fmt.Errorf("filePaths channel is nil")}
	}

	var (
		wg  = &sync.WaitGroup{}
		sem = &Semaphore{c: make(chan struct{}, fileDescriptorsLimit)}
		eP  = &ErrsPipe{errorsCh: make(chan error)}

		errS = make([]error, 1<<8)
	)

	for fp := range filePaths {

		wg.Add(1)
		go readFileС(wg, sem, eP, fp)

	}

	go func() {
		wg.Wait()
		close(eP.errorsCh)
		close(sem.c)
	}()

	for curErr := range eP.errorsCh {
		errS = append(errS, curErr)
	}

	return errS
}

func readFileС(wg *sync.WaitGroup, sem *Semaphore, eP *ErrsPipe, filePath string) {
	sem.c <- struct{}{}
	defer func() {
		<-sem.c
		wg.Done()
	}()

	f, err := os.Open(filePath)
	if err != nil {
		eP.errorsCh <- err
		return
	}
	defer func(filePath string) {
		fileClosingErr := f.Close()

		if fileClosingErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; closing file %q: %w", err, filePath, fileClosingErr)
			return
		}

		err = fmt.Errorf("closing file %q: %w", filePath, err)

		eP.errorsCh <- err
	}(filePath)

	// some logic...
}

func getProducerA() <-chan int {
	var (
		prod = make(chan int)
	)

	go func() {
		defer close(prod)

		for i := 1; i <= 100; i++ {

			prod <- i
		}

	}()

	return prod
}

func getProducerB() (<-chan int, <-chan struct{}) {
	var (
		prod = make(chan int, 100)
		done = make(chan struct{})
	)

	go func() {
		defer close(prod)

		for i := 1; i <= 100; i++ {
			prod <- i
		}

		close(done)
	}()

	return prod, done
}

func workWithProdA() {
	const (
		topBound = 50
	)

	var (
		prod = getProducerA()
	)

	for v := range prod {
		fmt.Println(v)

		if v == topBound {
			prod = make(<-chan int)
			continue
		}
	}
}

func workWithProdB() {
	defer fmt.Println("return from func")

	const (
		topBound = 50
	)

	var (
		prod, done = getProducerB()
	)
	<-done

	for {
		select {
		case v, ok := <-prod:
			if !ok {
				return
			}

			fmt.Println(v)
			if v == topBound {
				prod = make(<-chan int)
			}

		default:
			fmt.Println("new producer provided")
			return
		}
	}

}
