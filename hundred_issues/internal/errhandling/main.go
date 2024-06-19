package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

/*
	ERROR HANDLING

	PANIC
1. panic works in the following way:
	1) It resumes up on the stack until it'll be handled by recover function | there will be exit from the current goroutine (G)
2. recover() function is useful only in defer section. If there was no panic situation, it'll return nil value and nothing happens
3. defer() section is executed even if there's panic situation
4. panic is used only
	1) When there's a relevant exceptional critical situation (for example: an error of a programmer)
	2) When an application requires a specific dependency, but it cannot be initiated at all. For example: function Compile and MustCompile in regexp pkg, the second
	one requires the compilation of a specific regular expression and if it cannot be compiled, it throws a panic.
There are cases when it's needed to make both things: to add context and to marker an error as a specific one.

	WRAPPING AN ERROR
1. Error wrapping - it's packing an error inside a wrapper-container, that makes an internal error accessible
2. There are 2 use-cases of error-wrapping
	1) Additional context adding
	2) Markering an error as a specific one
In both cases a caller can unwrap an error an deconstruct it, using all contexts and specific information about it
3. How to handle errors in general?
	1) Return an error directly. In some cases we cannot add an additional context and mark the error as a specific one. In this situation we just return the source err.
	2) Wrap error with the "%w" directive. It saves an internal error and add text information. A client can unwrap the error and work with specification/markering
	3) Using "%v" directive. This directive doesn't wrap an error, instead of doing this way, it turns an error into a different one. After doing it, the source error starts
	to be unreachable. The caller cannot unwrap error. The info is acceptable, but the internal error is not.
4. Summary:
	1) Wrapping - is addition context of an error that makes the error more elaborative and verbose and/or markering it as a specific one.
	2) If we need to mark an error, we must create a specific type to do it.
	3) If we need just add an additional context to an error, we should just use the "%w" directive
	4) Nonetheless, during wrapping a potential relation created, since an internal error gets accessible to a caller. If we need to prevent this behavior, we should use the
	"%v" directive. "%v" directive doesn't wrap an error, it creates another error, but saves the specification of it.

	INACCURATE ERROR CHECK
1. Using custom error type | "%w" directive we also need to change error check of a caller code.
2. We should differentiate errors by their types. For this purpose we can use switch statement with .(type) to get the type of the error gotten.
3. If we've wrapped an error on the low-level and want to find out whether this error has the specific internal error type, we should use the errors.As(...) function.
4. erros.As(...) goes through the error structure recursively and checks whether the given error has the specific internal type. If there's an error that corresponds to the
type of error give, it'll return true.
	1) errors.As(...) requires that a target must be a pointer. If the target isn't a pointer, it'll panic.
5. If we use error wrapping we should use errors.As(...) to check whether the error has the specific type given.

*/

func main() {
	// PanicA()

	ErrorWrappingUsageA()
}

func PanicA() {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("recover in PanicA has handled", v.(string))
		}
	}()
	PanicB()
}

func PanicB() {
	PanicC()
}

func PanicC() {
	panic("panic in PanicC ")
}

func ErrorWrappingUsageA() {
	srcErr := errors.New("source error")
	wrappedError := wrap(srcErr, WrapArgs{Mark: "specific"})

	fmt.Println(wrappedError)
}

func wrap(err error, ewa WrapArgs) error {

	var (
		msg string
	)

	if len(ewa.Specs) != 0 {
		msg += fmt.Sprintf("%s: ", ewa.Specs)
	}

	if len(ewa.Mark) != 0 {
		switch l := len(msg); {
		case l != 0:
			msg = fmt.Sprintf("%s, ", ewa.Mark) + msg
		default:
			msg = fmt.Sprintf("%s: ", ewa.Mark)
		}
	}

	msg += "%w"

	return fmt.Errorf(msg, err)
}

func ErrHandling(err error) error {
	switch rand.Intn(4) {
	// No additional describing context and a need to mark an error as a specific
	case 0:
		return err
		// There's only an additional mark of an err
	case 1:
		return wrap(err, WrapArgs{Mark: "specific"})
		// There's only an additional specification of an err
	case 2:
		return wrap(err, WrapArgs{Specs: "while exec"})
		// There are both a specification and a mark of an err
	case 3:
		return wrap(err, WrapArgs{Mark: "specific", Specs: "while exec"})
	default:
		// Using the "%v" directive, turning an error into a different one
		return fmt.Errorf("while exec: %v", err)
	}
}

type WrapArgs struct {
	Mark  string
	Specs string
}

type CustomError struct {
	Err   error
	Msg   string
	Miscs map[string]any
}

func (ce CustomError) Error() string {
	return ce.Msg
}

func Wrap(err error, ewa WrapArgs) CustomError {

	var (
		msg string
	)

	if len(ewa.Specs) != 0 {
		msg += fmt.Sprintf("%s: ", ewa.Specs)
	}

	if len(ewa.Mark) != 0 {
		switch l := len(msg); {
		case l != 0:
			msg = fmt.Sprintf("%s, ", ewa.Mark) + msg
		default:
			msg = fmt.Sprintf("%s: ", ewa.Mark)
		}
	}

	msg += "%w"

	return CustomError{Err: err, Msg: msg, Miscs: make(map[string]any, 1<<4)}
}

func handler(w http.ResponseWriter, r *http.Request) {
	txID := r.URL.Query().Get("transaction")
	amount, err := getTxAmount(txID)
	if err != nil {
		// In this case the errors.As(err error, target any) function goes through the err chain given and matches all the errors encountered to the target type.
		// If there are any that match our target, it'll return true and set the target to the value
		if errors.As(err, &TransientError{}) {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// a text of the response

	fmt.Println(amount)
}

type TransientError struct {
	err error
}

func (t TransientError) Error() string {
	return fmt.Sprintf("transient error: %v", t.err)
}

func getTxAmount(txID string) (float32, error) {
	const (
		rightTxIDLen = 5
	)

	if len(txID) != rightTxIDLen {
		// Returning a simple err
		return 0, fmt.Errorf("id is invalid: %s", txID)
	}

	amount, err := getTxAmount(txID)
	if err != nil {
		// Returning a custom err wrapping the transient error come from getTxAmountFromDB
		return 0, fmt.Errorf("failed to get tx %s: %w", txID, err)
	}

	return amount, nil
}

func getTxAmountFromDB(txID string) (float32, error) {
	var (
		amount float32
	)

	for _, r := range txID {
		amount += float32(rand.Intn(int(r)))
	}

	if rand.Intn(2) == 1 {
		return amount, TransientError{err: errors.New("an error")}
	}

	return amount, nil
}
