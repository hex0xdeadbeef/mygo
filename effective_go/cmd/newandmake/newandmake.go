package newandmake

import (
	"bytes"
	"fmt"
	"sync"
)

type (
	SyncedBuffer struct {
		lock   sync.Mutex
		buffer bytes.Buffer
	}
)

func zeroedValueUsing() {
	p := new(SyncedBuffer)
	v := SyncedBuffer{}

	fmt.Println(p)
	fmt.Println(v)
}
