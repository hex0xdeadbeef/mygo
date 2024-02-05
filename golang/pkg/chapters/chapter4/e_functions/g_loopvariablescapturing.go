package e_functions

import (
	"fmt"
	"os"
)

var (
	rmdirs func()
)

func WorkWithAndRemoveDirectories() {
	requiredNumberOfDirs := 10

	removeDirs := make([]func(), 0, requiredNumberOfDirs)
	// Create some temporary directories and append functions to remove them
	for _, curDir := range tempDirs(requiredNumberOfDirs) {
		localLoopDir := curDir
		os.MkdirAll(localLoopDir, 0755)

		/*
			Every closure refers to local mempool of curDir variable, so finally
			only the last directory will be removed
		*/
		removeDirs = append(removeDirs, func() { os.RemoveAll(localLoopDir) })
	}

	/*
		... do some work ...
	*/

	for _, removeFunc := range removeDirs {
		removeFunc()
	}
}

func tempDirs(tempDirsNumber int) []string {
	tempDirs := make([]string, tempDirsNumber)
	for i := 0; i < tempDirsNumber; i++ {
		tempDirs = append(tempDirs, os.TempDir())
	}

	return tempDirs
}

func CapturingSliceElementByIndex() []func() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	funcs := make([]func(), 0, 10)

	for i := 0; i < len(numbers); i++ {
		index := i
		funcs = append(funcs, func() { fmt.Println(&i, &index, numbers[index]) })

	}

	return funcs
}
