package readerdonewrapper

/*
OrDone channel allows us to check whether the done channel has been already closed, at the moment we read from it.
*/

func IsReaderDone(done <-chan struct{}, dataProducer <-chan interface{}) <-chan interface{} {
	// Create the additional channel we will send data into
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)

		for {
			select {
			// If the done channel has been already closed
			case <-done:
				return
			// Check whether the dataProducer has been already closed
			case val, ok := <-dataProducer:
				if !ok {
					return
				}

				// Repeatedly check whether the done channel is closed or not at the moment we've reached this statements
				select {
				case <-done:
					return
					// If not, send the given value
				case valStream <- val:
				}
			}
		}
	}()

	return valStream
}
