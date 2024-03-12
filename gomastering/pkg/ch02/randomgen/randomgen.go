package randomgen

import (
	"bytes"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	seed int64
	src  rand.Source
	rnd  *rand.Rand
)

func init() {
	os.Args = append(os.Args, "0 127 100")

	seed = time.Now().Unix()
	src = rand.NewSource(seed)
	rnd = rand.New(src)

	log.Print(seed)
}

type RandInts struct {
	Min, Max int
	Count    int
	Rands    []int
}

func (ri RandInts) Print() {
	const maxNumberInLine = 10

	for i := 0; i < ri.Count; i++ {
		fmt.Printf(" %d", ri.Rands[i])

		if i != 0 && i%maxNumberInLine == 0 {
			fmt.Println()
		}
	}
}

func RandNumbers() (RandInts, error) {
	args := strings.Fields(os.Args[1])

	if len(args) != 3 {
		return RandInts{}, fmt.Errorf("Usage: min(int) max(int) count(positive int)")
	}

	ri, err := setParams(args)
	if err != nil {
		return RandInts{}, fmt.Errorf("validating params: %w", err)
	}

	ri.Rands = make([]int, ri.Count)

	fill(ri)

	return ri, nil
}

func setParams(args []string) (RandInts, error) {
	if len(args) != 3 {
		return RandInts{}, fmt.Errorf("invalid params count: %d", len(args))
	}

	ri := RandInts{}
	reflInts := reflect.ValueOf(&ri).Elem()

	for i, arg := range args {
		if n, err := strconv.Atoi(arg); err != nil {
			return RandInts{}, nil
		} else {
			reflInts.Field(i).SetInt(int64(n))
		}
	}

	return ri, nil
}

func fill(ri RandInts) {
	for i := 0; i < ri.Count; i++ {
		ri.Rands[i] = rnd.Intn(ri.Max-ri.Min) + ri.Min
	}
}

func GetPass() (string, error) {
	const (
		initialBufferSize = 2048
	)

	buf := bytes.NewBuffer(make([]byte, initialBufferSize))

	ri, err := RandNumbers()
	if err != nil {
		return "", fmt.Errorf("generating numbers %s", err)
	}

	for _, val := range ri.Rands {
		if !unicode.IsPrint(rune(val)) {
			continue
		}

		if _, err := buf.Write([]byte(string(rune(val)))); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func SecureRandPass(n int64) (string, error) {
	b := make([]byte, n)

	if _, err := crand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
