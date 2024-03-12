package which

import (
	"fmt"
	"os"
	"path/filepath"
)

func Which() error {
	const (
		invalidArgsCount = 1

		envVariableName = "PATH"
	)
	args := os.Args

	if len(args) == invalidArgsCount {
		return fmt.Errorf("checking args count: invalid args count: %d", len(args))
	}

	for _, fileName := range os.Args[1:] {
		path := os.Getenv(envVariableName)
		if path == "" {
			return fmt.Errorf("invalid environment variable's name provided: %s", envVariableName)
		}
		pathSplit := filepath.SplitList(path)

		for _, directory := range pathSplit {
			fullPath := filepath.Join(directory, fileName)

			if fileInfo, err := os.Stat(fullPath); err == nil {
				if mode := fileInfo.Mode(); mode.IsRegular() && mode&0111 != 0 {
					fmt.Println(fullPath)
				}
			}
		}
	}

	return nil
}
