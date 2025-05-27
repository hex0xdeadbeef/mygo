package main

import "fmt"

func WrapDeferredErr(accumulatedErr error, ft string, f func() error) error {
	resErr := f()

	if accumulatedErr == nil && resErr == nil {
		return nil
	}

	if accumulatedErr == nil && resErr != nil {
		return fmt.Errorf(ft, resErr)
	}

	if resErr == nil && accumulatedErr != nil {
		return fmt.Errorf(ft, accumulatedErr)
	}

	resErr = fmt.Errorf(ft, resErr)

	return fmt.Errorf("%w; %w", accumulatedErr, resErr)
}
