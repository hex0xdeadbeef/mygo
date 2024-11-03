package main

func main() {
	var m = map[string]struct{}{"": {}}

	go func() {
		for k, v := range m {
			_, _ = k, v
		}
	}()

	for k, v := range m {
		_, _ = k, v
	}
}
