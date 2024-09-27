package chapter7

import (
	"database/sql"
	"fmt"
	"log"
)

func listTracks(db sql.DB, artist string, minYear, maxYear int) {
	result, err := db.Exec(
		"SELECT * FROM tracks WHERE artist = ? AND ? <= year <= ?", artist, minYear, maxYear)
	if err != nil {
		log.Fatal(err)
	}

	// ... do some work

	fmt.Println(result)

}

func sqlQuoteIfs(x interface{}) string {
	if x == nil {
		return "NULL"
	} else if _, ok := x.(int); ok {
		return fmt.Sprintf("%d", x)
	} else if _, ok := x.(uint); ok {
		return fmt.Sprintf("%d", x)
	} else if b, ok := x.(bool); ok {
		if b {
			return "TRUE"
		}
		return "FALSE"
	} else if s, ok := x.(string); ok {
		return sqlQuoteString(s)
	} else {
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}

func SqlQuoteSwitch(x interface{}) string {
	switch x.(type) {
	case nil:
		return "NULL"
	case int, uint:
		return fmt.Sprintf("%d", x)
	case bool:
		x := x.(bool)
		if x {
			return "TRUE"
		}
		return "FALSE"
	case string:
		x := x.(string)
		return sqlQuoteString(x)
	default:
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}

func SqlQuoteSwitchExtended(x interface{}) string {
	// There's the new x variable
	switch x := x.(type) {
	case nil:
		fmt.Printf("%T\n", x)
		return "NULL"
	case int, uint:
		fmt.Printf("%T\n", x)
		return fmt.Sprintf("%d", x) // x has the interface{} type
	case bool:
		fmt.Printf("%T\n", x)
		if x {
			return "TRUE"
		}
		return "FALSE"
	case string:
		fmt.Printf("%T\n", x)
		return sqlQuoteString(x)
	default:
		fmt.Printf("%T\n", x)
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}

func sqlQuoteString(s string) string {
	// ... do some work
	return s
}
