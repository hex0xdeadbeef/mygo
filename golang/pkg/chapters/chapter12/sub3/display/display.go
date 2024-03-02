package display

import (
	"fmt"
	"golang/pkg/chapters/chapter12/sub2/format"
	"reflect"
)

func DisplayUsing() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actors          map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actors: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
		},
	}

	Display("strangelove", strangelove)
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))

}

type empty struct{}

var (
	repetitions map[reflect.Value]empty
)

func init() {
	repetitions = make(map[reflect.Value]empty)
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice {
			if isCycle(v) {
				return
			}

			repetitions[v] = empty{}
		}

		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		if isCycle(v) {
			return
		}

		repetitions[v] = empty{}

		if keys := v.MapKeys(); len(keys) != 0 {
			firstKey := keys[0]

			// If a internal value isn't comparable we'll break the case statement
			switch firstKey.Kind() {
			case reflect.Array:
				for _, key := range keys {
					strKey := fmt.Sprintf("%s[%v]", path, key)
					display(fmt.Sprintf("%s[%s]", path, strKey), v.MapIndex(key))
				}
			case reflect.Struct:
				for _, key := range keys {
					strKey := fmt.Sprintf("%s[%v]", path, key)
					display(fmt.Sprintf("%s[%s]", path, strKey), v.MapIndex(key))
				}
			default:
				for _, key := range keys {
					display(fmt.Sprintf("%s[%s]", path, format.FormatAtom(key)), v.MapIndex(key))
				}
			}
			return
		}

		fmt.Printf("map has 0 length")

	case reflect.Ptr:
		if isCycle(v) {
			return
		}

		repetitions[v] = empty{}

		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
			return
		}
		display(fmt.Sprintf("(*%s)", path), v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
			return
		}
		fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
		display(path+".value", v.Elem())
	default:
		fmt.Printf("%s = %s\n", path, format.FormatAtom(v))
	}

}

func isCycle(val reflect.Value) bool {
	if _, ok := repetitions[val]; ok {
		fmt.Printf("CYCLE\n")
		delete(repetitions, val)
		return true
	}
	return false
}

func DifferencePtrAndNonPtrDisplaying() {
	var (
		number                           = 3
		emptyInterfaceNumber interface{} = number
	)

	Display("number", number)
	fmt.Println()
	Display("&number", &number)
	fmt.Println()
	Display("&emptyInterfaceNumber", &emptyInterfaceNumber)
}
