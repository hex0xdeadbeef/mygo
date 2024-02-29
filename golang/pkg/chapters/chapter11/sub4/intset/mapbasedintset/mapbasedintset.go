package mapbasedintset

import "fmt"

type empty struct{}

type MapIntSet map[uint]empty

func New() *MapIntSet {
	underlyingMap := make(map[uint]empty)
	set := MapIntSet(underlyingMap)

	return &set
}

func (mis *MapIntSet) Add(values ...int) {
	for _, value := range values {
		(*mis)[uint(value)] = empty{}
	}
}

func (mis *MapIntSet) Clear() error {
	if mis != nil {
		(*mis) = make(map[uint]empty)
		return nil
	}

	return fmt.Errorf("map int set ptr is nil")
}
