package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

type empty struct{}

var (
	repetitions map[reflect.Value]empty
)

func init() {
	repetitions = make(map[reflect.Value]empty)
}

func isPresented(val reflect.Value) bool {
	if _, ok := repetitions[val]; ok {
		return true
	}
	return false
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	var (
		kind = v.Kind()
	)

	if kind != reflect.Invalid && v.IsZero() {
		return nil
	}

	switch kind {
	// Invalid
	case reflect.Invalid:
		buf.WriteString("nil")
		// Ints
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
		// Uints
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
		// Float
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%.5f", v.Float())
		// Complex
	case reflect.Complex64, reflect.Complex128:
		compl := v.Complex()
		fmt.Fprintf(buf, "#C (%.5f %.5f)", real(compl), imag(compl))
	case reflect.String:
		// String
		fmt.Fprintf(buf, "%q", v.String())
		// Bool
	case reflect.Bool:
		val := v.Bool()
		if val {
			fmt.Fprint(buf, "t")
			break
		}
		fmt.Fprintf(buf, "nil")
		// Pointer
	case reflect.Ptr:
		if isPresented(v) {
			delete(repetitions, v)
			return nil
		}
		repetitions[v] = empty{}
		return encode(buf, v.Elem())
	case reflect.Interface:
		if isPresented(v) {
			delete(repetitions, v)
			return nil
		}
		repetitions[v] = empty{}

		if v.IsNil() {
			return fmt.Errorf("nil interface value")
		}
		buf.WriteByte('(')
		fmt.Fprintf(buf, "%q ", v.Type().Name())
		if err := encode(buf, v.Elem()); err != nil {
			return err
		}
		buf.WriteByte(')')
		// Array/Slice (value ...)
	case reflect.Array, reflect.Slice:
		if v.Kind() == reflect.Slice {
			if isPresented(v) {
				delete(repetitions, v)
				return nil
			}
			repetitions[v] = empty{}
		}

		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
		// Struct ((name value) ...)
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')

		}
		buf.WriteByte(')')

		// Map ((key value) ...)
	case reflect.Map:

		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')

			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}

			buf.WriteByte(')')

		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type().String())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Print(p []byte) {
	const (
		scale                     = 2
		doubleOpenBrackets        = "(("
		intendeddifferentBrackets = ")\n"
		doubleCloseBrackets       = "))"

		closeBracket = ")"
		openBracket  = "("
	)
	var (
		pair string

		depth = 0
	)

	p = bytes.ReplaceAll(p, []byte(") ("), []byte(")\n("))

	for i := 0; i < len(p)-1; i++ {
		pair = string(p[i]) + string(p[i+1])

		switch pair {
		case doubleOpenBrackets:
			fmt.Print(doubleOpenBrackets)
			depth++
			i++
		case intendeddifferentBrackets:
			fmt.Println(closeBracket)
			fmt.Printf("%*s", scale*depth, "")
			i++
		case doubleCloseBrackets:
			fmt.Printf("%s\n", doubleCloseBrackets)
			depth--
			i++
		default:
			fmt.Print(string(p[i]))
		}
	}
}

type lexer struct {
	scan scanner.Scanner
	// the current token
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	// NOTE: Not an example of good error handling
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}

	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	const (
		nilStr = "nil"
	)
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are "nil"
		// and struct field names.
		if lex.text() == nilStr {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		// NOTE: ignoring errors
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return

	case scanner.Int:
		// NOTE: ignoring errors
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		// consume ')'
		lex.next()
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	// (item ...)
	case reflect.Array:
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}
	// (item ...)
	case reflect.Slice:
		var item reflect.Value

		for !endList(lex) {
			item = reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}
	// ((name value) ...)
	case reflect.Struct:
		var name string
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name = lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}
	// ((key value) ...)
	case reflect.Map:
		var (
			key   reflect.Value
			value reflect.Value
		)

		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
		}
		key = reflect.New(v.Type().Key()).Elem()
		read(lex, key)

		value = reflect.New(v.Type().Elem()).Elem()
		read(lex, value)
		v.SetMapIndex(key, value)
		lex.consume(')')
	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("EOF")
	case ')':
		return true
	}
	return false
}

// Unmarshall parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshall(data []byte, out interface{}) (err error) {
	var lex *lexer
	lex = &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	// get the first token
	lex.next()

	// NOTE: this isn't an example of ideal error handling
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()

	read(lex, reflect.ValueOf(out).Elem())
	return nil
}
