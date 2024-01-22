package b_count_of_terminal_arguments

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Dup_First() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(file *os.File, counts map[string]int) {
	input := bufio.NewScanner(file)
	for input.Scan() {
		counts[input.Text()]++
	}
}

func Dup_Second(filesNames ...string) {
	counts := make(map[string]int)

	if len(filesNames) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, fileName := range filesNames {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Dup: %v\n", err)
				continue
			} else {
				countLines(file, counts)
				file.Close()
			}
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func Dup_Third() {
	counts := make(map[string]int)

	for _, filename := range os.Args[1:] {
		dataFromFile, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprint(os.Stderr, "Dup_Third:\n", err)
			continue
		}

		for _, line := range strings.Split(string(dataFromFile), "\n") {
			counts[line]++
		}

	}

	for line, number := range counts {
		if number > 1 {
			fmt.Printf("%d : %s", number, line)
		}
	}
}

// 1.3 Homework
func Dup_Second_Modified(filesNames ...string) {

	if len(filesNames) == 0 {
		counts := make(map[string]int)

		countLines(os.Stdin, counts)

		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	} else {
		fileCounts := make(map[int]map[string]int) // Create map of maps

		// Fill a map for each file
		for i, fileName := range filesNames {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error ocurred: file doen't exist")
				continue
			} else {
				counts := make(map[string]int) // Creates an empty map
				countLines(file, counts)       // Pass a map and file into countLines() to be filled
				fileCounts[i] = counts         // Add a filled map into map of maps
				file.Close()                   // Close a file
			}
		}

		unionedMap := mapUnion(fileCounts) // Create an union map

		// Go through the unioned map, print count, string, files names strings are contained in

		for key, value := range unionedMap {
			fmt.Printf("|%d|\t|%s|\t", value, key)
			for i, curMap := range fileCounts {
				if _, ok := curMap[key]; ok {
					fmt.Printf("|%s|\t", filesNames[i])
				}
			}
			fmt.Println()
		}

	}

}

func mapUnion(mapOfMaps map[int]map[string]int) map[string]int {
	unionedMap := make(map[string]int)
	for _, curMap := range mapOfMaps {
		for key, value := range curMap {
			if _, ok := unionedMap[key]; !ok {
				unionedMap[key] = value
			} else {
				unionedMap[key] += value
			}

		}
	}

	fmt.Println()
	return unionedMap
}
