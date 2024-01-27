package d_3json

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
1. JSON (Java Script Object Notation) - is a standart notation for sending and receiving structured information with HTTP(s)
protocol(s)

2. The basic JSON types are:
	- Numbers (decimal or scientific notation)

	- Booleans (true, false)

	- Strings (sequences of Unicode code points enclosed in double quotes)
		! JSON'S \Uhhhh NUMERIC ESCAPES DENOTE UTF-16 CODES !

3. Recursively combined basic types:
	- Array [value, value, value, ... ] is used to encode Go's arrays/slices

	- Object {str1:value, ... strN:value } is a mapping: strings -> values separated by commas and surrounded by figure braces.
	Object is used to represent Go's structs/maps(with string keys).

4. Converting Go's objects to its JSON representation is called "marshaling".

5. json.Marshal(data) produces:
	1) []byte containing a very long string

	2) error if something goes wrong
	! THE RESULT OF THIS OPERATION IS IMPOSSIBLE TO BE READ !

6. json.MarshalIndent(data, "prefix", "indentation") produces neatly indented output

7. Marshaling uses Go's struct fields names as the field names for the JSON objects (Through reflection).

8. Only exported fields are marshaled

9. Field tag is the mechanism for creating metadata for original data:
	It's declared as: `key:"value1, ..., valueN"`
	1) The first parameter produces another JSON type name so that struct filed name fits the json field name

	2) "omitempty" indicates that no JSON output should be produced if the field has zero value for its type or is otherwise
	empty.

10. The inverse operation is called Unmarshaling, it's done by json.Unmarshal

11. Unmarshaling operation is case-insensitive, but when JSON has underscored fields it's only necessary to provide the cor-
responding field with a tag.
*/

type Movie struct {
	Title  string
	Year   int  `json:"released"`        // It's a field tag
	Color  bool `json:"color,omitempty"` // It's a field tag
	Actors []string
}

var (
	movies = []Movie{
		{Title: "Casablanca", Year: 1942, Color: false, Actors: []string{"Humpey Bogart", "Ingrid Bergman"}},
		{Title: "Cool Hand Luke", Year: 1967, Color: true, Actors: []string{"Paul Newman"}},
		{Title: "Bullit", Year: 1968, Color: true, Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	}
)

func EncodeMarshal() {

	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("JSON marshaling failed %s", err)
	}

	fmt.Printf("data after marshaling:\n%v", string(data))
	/*[
	{"Title":"Casablanca","released":1942,"Color":false,"Actors":["Humpey Bogart","Ingrid Bergman"]},
	{"Title":"Cool Hand Luke","released":1967,"Color":true,"Actors":["Paul Newman"]},
	{"Title":"Bullit","released":1968,"Color":true,"Actors":["Steve McQueen","Jacqueline Bisset"]}
	]
	*/
}

func EncodeMarshalIndent() {
	data, err := json.MarshalIndent(movies, "", "\t")
	if err != nil {
		log.Fatalf("JSON indented marshaling failed: %s", err)
	}

	fmt.Printf("data after indented marshaling:\n%v", string(data))
}

func DecodeUnmarshal() {
	data, err := json.MarshalIndent(movies, "", "\t")
	if err != nil {
		log.Fatalf("While marshaling error occured: %s", data)
	}

	var titles []struct{ Title string } // Create an anonymous structure in-place
	if err := json.Unmarshal(data, &titles); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	fmt.Println(titles)
}
