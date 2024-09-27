package chapter10

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

type empty struct{}

// PackageNode is used to store data after parsing the import paths
type PackageNode struct {
	ImportPath string
	Imports    []string
}

// TransitiveSet is used to save the identity of root while traversing all the imports
type TransitiveSet struct {
	root                *PackageNode
	orderedDependencies []string

	seen  map[string]empty
	queue []string
}

// NewTransitiveSet creates the new *TransitiveSet instance based on pn
func NewTransitiveSet(pn *PackageNode) *TransitiveSet {
	ts := &TransitiveSet{root: pn, seen: make(map[string]empty), queue: make([]string, 0, len(pn.Imports))}

	ts.seen[pn.ImportPath] = empty{}
	ts.queue = append(ts.queue, ts.root.ImportPath)

	return ts
}

// FindTransitivePackages is the controller function that combines all the logic, traverse all the roots given and finds the transitive import chains
// returns an error encountered if any
func FindTransitivePackages() error {
	var (
		sets []*TransitiveSet
	)

	pkgsRoots, err := roots()
	if err != nil {
		return fmt.Errorf("getting roots from args: %s", err)
	}
	sets = make([]*TransitiveSet, 0, len(pkgsRoots))

	for _, pkgRoot := range pkgsRoots {
		wg.Add(1)
		go func(pkgRoot string) {
			defer wg.Done()
			setRoot(pkgRoot, &sets)
		}(pkgRoot)
	}

	wg.Wait()

	for _, ts := range sets {
		wg.Add(1)
		go func(ts *TransitiveSet) {
			defer wg.Done()
			ts.traverseImports()
		}(ts)
	}

	wg.Wait()

	for _, set := range sets {
		set.printTransitiveDependencies()
	}

	return nil
}

// roots checks validities of prompt's args and if everything's okay returns args and an error encountered if any
func roots() ([]string, error) {
	const (
		minPromptArgsCount = 1
	)

	if len(os.Args[1:]) < minPromptArgsCount {
		return nil, fmt.Errorf("insufficient args count: %d", len(os.Args[1:]))
	}

	for _, pkgPath := range os.Args[1:] {
		if !filepath.IsLocal(pkgPath) {
			return nil, fmt.Errorf("parsing paths: \"%s\" is not a path", pkgPath)
		}
	}

	return os.Args[1:], nil
}

// goListJSON gets the result of go list -json with provided path and returns parsed object and an error encountered if any
func goListJSON(pkgPath string) (*PackageNode, error) {
	cmd := exec.Command("go", "list", "-json", pkgPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("executing go list -json \"%s\": %s", pkgPath, err)
	}

	var pkgNode PackageNode
	err = json.Unmarshal(output, &pkgNode)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %s", err)
	}

	return &pkgNode, nil
}

// setRoot adds the new *TransitiveSet instance into the slice of objects to be recursively traversed
func setRoot(path string, sets *[]*TransitiveSet) {
	root, err := goListJSON(path)
	if err != nil {
		log.Printf("parsing \"%s\": %s", path, err)
		return
	}

	mu.Lock()
	*sets = append(*sets, NewTransitiveSet(root))
	mu.Unlock()
}

// traverseImports traverses all the imports of initial root package
func (ts *TransitiveSet) traverseImports() {
	var (
		queueHead       string
		nextPackageNode *PackageNode = ts.root

		err error
	)

	for len(ts.queue) != 0 {

		queueHead = ts.queue[0]
		nextPackageNode, err = goListJSON(queueHead)
		if err != nil {
			log.Print(err)
			continue
		}

		if len(ts.queue) == 1 {
			ts.queue = nil
		} else {
			ts.queue = ts.queue[1:]
		}

		for _, importPath := range nextPackageNode.Imports {
			_, ok := ts.seen[importPath]
			if !ok {
				ts.seen[importPath] = empty{}
				ts.orderedDependencies = append(ts.orderedDependencies, importPath)

				nextImports, err := goListJSON(importPath)
				if err != nil {
					log.Print(err)
					return
				}

				ts.queue = append(ts.queue, nextImports.Imports...)

				// mark children of the current node as discovered and add them into orderedDependencies
				for _, newImportPath := range nextImports.Imports {
					ts.seen[newImportPath] = empty{}
					ts.orderedDependencies = append(ts.orderedDependencies, newImportPath)
				}

			}
		}

	}
}

// printTransitiveDependencies prints the ordered transitive dependencies in a table format
func (ts *TransitiveSet) printTransitiveDependencies() {
	const (
		header = "Root: %s\n"
		row    = "%*s\n"
	)

	var (
		rowIndent = len(header) + len(filepath.Base(ts.root.ImportPath))
	)

	fmt.Printf(header, filepath.Base(ts.root.ImportPath))
	for _, dependency := range ts.orderedDependencies {
		fmt.Printf(row, len(dependency)+rowIndent, dependency)
	}
}
