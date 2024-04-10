package math

// Finds an average value in floats slice
func Average(numbers []float64) float64 {
	var average float64 = 0
	for _, element := range numbers {
		average += element
	}
	return average / float64(len(numbers))
}

// Finds a min value in floats slice
func Min(numbers []float64) float64 {
	var min float64 = 0
	for _, element := range numbers {
		if element < min {
			element = min
		}
	}
	return min
}

// Finds a max value in floats slice
func Max(numbers []float64) float64 {
	var max float64 = 0
	for _, element := range numbers {
		if element > max {
			element = max
		}
	}
	return max
}
