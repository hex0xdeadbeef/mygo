package c_hashmaps

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"unicode"
	"unicode/utf8"
)

/*
1. Map is an unordered collection of key/value pairs
2. All keys are distinct from each other
3. The value associated with a given key can be retrieved, updated, removed, using CONSTANT NUMBER OF KEY COMPARISONS
no matter how large hash table.
4. Map type is map[k]v
5. The k type must be comparable using ==
6. The v type may have any type
7. Accessing, changing, deleting a value
	1) We can access the value with its key. If map doesn't consist the key:
		1) Using mapName[key] = value, we ADD the new key/value pair if it doesn't exist in a map
		2) Using value := mapName[key], we get the current value corresponding to the key, if key isn't contained in
		map, the ZERO value of its type and BOOL value
	2) We can change the value using a key
	3) We can delete a value using operator (from, key)
8. We can't take the address of the map element because of elements migration. Elements migration is the state when
an average amount of a map values reachs 6.5 and data is moved into new storage respectively. Another cause of this
event is rehasing elements as the number of elements growth, thus potentially invalidating the addresses.
9. The order of enumerating maps is arbitrary because of difference between implementations of map. In implementatio
ns different hash funcions are used.
10. So that enumerate the pairs key/value of a map orderly we must explicitly sort the keys. For efficiency we should
allocate a slice with a capacity required to store all the keys in advance (its length is map length).
11. For better efficiency we should allocate memory for slice storing the keys of a map in advance
12. The zero value of map is nil, that is, a reference to no hash table at all
13. Map doesn't have capacity
14. + Allowed operations on a nil map: lookup, delete, len, for loops
	1) lookup will return 0
	2) delete always returns nothing after it's called, so after deleting an nonexisting pair, nothing fill happen
	3) The length of a nil map is 0
	4) enumeration with for loop will do nothing
15. - Forbidden operation on a nil map: storing (mapName[nonexisting key] = value) // panic
16. Before storing to a map, we should create it
17. Accessing a map value always yields a value even if the value isn't present in the map
	1) If the key is present in the map, we get the corresponding value
	2) If not we get the zero value for the element type
18. So that distinguish whether a value is nonexising or it happens to have the zero value we use the following sta-
tement:
	value, isExist := mapName[key] // depending whether key is present in map or not isExist will have values true/
	false
19. Maps can be compared only with nil because of the migration of the data it consists
20. Go doesn't provide a "set" type, but since keys of a map are distinct, a map can serve this purpose.
21. When we need a map/set to deal with keys of non-comparable types (slices, other maps and so on)
	1) We should create a function k(x) that maps for each argument of type "x"a string.
	2) Using a values from k(x) (strings) we create a map which keys are strings.
22. WE CAN DEFINE THE FUNCTION k(x) AS FUNCTION THAT RECEIVES ANY ARGUMENTS OF ANY TYPES, MODIFIES THEM AND RETURNS
	MODIFIED VALUES SO THAT THEY WILL MATCHED WITH ==. P.S. N.L. Dodonova krasava as fuck
	*** THIS WAY ALLOWS US TO DEFINE OUR LOGIC OF > < == != ***
23. If we try to get a value from nil map, we'll get the zero value of its type
*/

func MapCreation() {
	ages1 := make(map[string]int) // It's the prefilled map, we can store data
	ages2 := map[string]int{}     // It's the prefilled map, we can store data

	ages3 := map[string]int{"Dmitriy": 21, "Katya": 20} // It's the initially filled map with some values with
	// map literal
	ages4 := make(map[string]int) // It's the initially prefilled map that will be filled with values in
	// the following statements
	ages4["Dmitriy"] = 20
	ages4["Maxim"] = 21

	fmt.Printf("%v\n%v\n%v\n%v\n", ages1, ages2, ages3, ages4)

}

func GetAddRemoveValue() {

	ages := map[string]int{"Dmitriy": 21, "Katya": 20}
	fmt.Println("%v", ages)

	// Change a value with key using
	ages["Dmitriy"] = 100
	fmt.Println("%v", ages)

	// Add a new value with an nonexisting key within the map
	ages["Maxim"] = 200
	fmt.Println("%v", ages)

	// Deleting the value using delete() operator
	delete(ages, "Maxim")
}

func NonexistingKeyAccess() {
	// The initial map
	ages := map[string]int{"Dmitriy": 21, "Katya": 20}
	fmt.Printf("%v\n", ages) // [Dmitriy:21 Katya:20]

	// Accessing the value with an nonexisting key returns a type corresponding zero value
	zero := ages["Margarita"]
	fmt.Printf("Returned zero value %d\n", zero)

	// The map doesn't have the value of Margarita's age yet // map[Dmitriy:21 Katya:20]
	fmt.Printf("%v\n", ages)

	fmt.Printf("Dmitriy's age: %d, Margarita's age: %d\n", ages["Dmitriy"], ages["Margarita"])

	// Check the value of the nonexisting key
	zeroAgeAlexey := ages["Alexey"]
	// There is absence of the Alexey's age // 0
	fmt.Printf("Alexey's age is: %d\n", zeroAgeAlexey)

	// Add the Alexey's age using the returned zero value and incrementing it
	ages["Alexey"] += 1 // He's been born now
	// Now the map has the following values: map[Alexey:1 Dmitriy:21 Katya:20]
	fmt.Printf("%v\n", ages)

	// Check the value of the added key/value pair
	zeroAgeAlexey = ages["Alexey"]
	// There is absence of the Alexey's age // Alexey's age is: 1
	fmt.Printf("Alexey's age is: %d\n", zeroAgeAlexey)
}

func IncrementingValues() {
	// The initial map
	ages := map[string]int{"Dmitriy": 21, "Katya": 20}
	fmt.Printf("%v\n", ages) // [Dmitriy:21 Katya:20]

	// Add the Alexey's age using the returned zero value and incrementing it
	ages["Alexey"] += 1 // He's been born now
	// Now the map has the following values: map[Alexey:1 Dmitriy:21 Katya:20]
	fmt.Printf("%v\n", ages)

	// The folowing ones also work
	ages["Alexey"] = ages["Alexey"] + 1
	ages["Alexey"]++

	// The map now has the following pairs:
	fmt.Printf("%v\n", ages) // map[Alexey:3 Dmitriy:21 Katya:20]
}

/*
There will be the compiller error because when an average amount of map values gets 6.5, the data migration happens
*/
func TryToGetValueAddress() {
	// The initial map
	ages := map[string]int{"Dmitriy": 21, "Katya": 20}
	fmt.Printf("%v\n", ages) // [Dmitriy:21 Katya:20]

	// The attempt to get the address of value
	// pointer := &ages["Dmitriy"] // invalid operation: cannot take address of ages["Dmitriy"] (map index expression of type int)

}

func MapEnumerating() {
	// The initial map
	ages := map[string]int{"Dmitriy": 21, "Katya": 20, "Yulia": 40}
	// This for statement will print the map pairs without an order. The order will be arbitary.
	for key, value := range ages {
		fmt.Printf("key: %s\tvalue: %d\n", key, value)
	}
}

func OrderedEnumeration() {
	// The initial map
	ages := map[string]int{"Yulia": 40, "Dmitriy": 21, "Kate": 20}

	// Create a slice for the keys of the map
	// var keys []string

	/* It'll be better for efficiency if we explicitly create the slice of the sufficient capacity required to store
	all the keys */
	keys := make([]string, 0, len(ages)) // {}, len = 0, cap = 3

	// Add all the unordered keys to the slice
	for key, _ := range ages {
		keys = append(keys, key)
	}
	fmt.Printf("Keys: %v\n", keys)

	// Explicitly sort the slice of keys using the sort.Strings() function
	sort.Strings(keys)

	fmt.Printf("Sorted keys: %v\n\n", keys)

	// Ordered enumeration pairs key/value of the map
	for index, key := range keys {
		fmt.Printf("%d. Key: %10s\t->\tValue: %5d\n", index, key, ages[key])
	}
}

func ZeroValueOfMap() {
	var nilMap map[int]string
	fmt.Printf("The map is: %v | It's len: %d | Is map nil value?: %t\n", nilMap, len(nilMap), nilMap == nil)
}

func StoringNilMap() {
	var nilMap map[string]int

	/* Allowed operations on a nil map: lookup, len, delete, for loops */

	// lookup
	fmt.Printf("The value of the nonexisting key: %d\n", nilMap["Dmitriy"]) // 0
	// len()
	fmt.Printf("The length of the nil map is: %d\n", len(nilMap)) // 0
	// delete(source map, key)
	fmt.Printf("The map before deleting the pair with the key \"Dmitriy\": %v\n", nilMap) // []map
	delete(nilMap, "Dmitriy")
	fmt.Printf("The map after deleting the pair with the key \"Dmitriy\": %v\n", nilMap) // []map
	// enumeration with foor loop
	for index, value := range nilMap {
		fmt.Printf("index: %d \tvalue: %s\n", index, value)
	}

	/* Forbidden operation: storing (mapName[key] = value) // panic */
	nilMap["Name"] = 10
	/*
		Exception has occurred: panic
		"assignment to entry in nil map"
		Stack:
		3  0x00000001008e72c8 in golang/pkg/chapter4/c_hashmaps.StroingNilMap
			at /Users/dmitriymamykin/Desktop/goprs/golang/pkg/chapter4/c_hashmaps/c_hashmaps.go:197
		4  0x00000001008e730c in main.main
			at /Users/dmitriymamykin/Desktop/goprs/golang/cmd/main.go:8
	*/
}

func ExistenceCheck() {
	filledMap := map[string]int{"A": 0, "B": 2, "C": 3}

	// This set of statements is restricted
	if value, doesExist := filledMap["A"]; doesExist {
		fmt.Printf("The corresponding value to the key \"A\" is: %d\n", value) // IT'LL BE PRINTED
	} else {
		fmt.Printf("The value for the key \"A\" isn't present in the map.")
	}

	// This set of instructions isn't restricted, so values: value, doesExist can be used further
	value, doesExist := filledMap["D"]
	if doesExist {
		fmt.Printf("The corresponding value to the key \"A\" is: %d\n", value)
	} else {
		fmt.Printf("The value for the key \"D\" isn't present in the map.") // IT'LL BE PRINTED
	}

}

// HANDMADE MAPS COMPARISON
func MatchMapsData() {
	map1 := map[string]int{"A": 1, "B": 2, "C": 3}
	map2 := map[string]int{"A": 0, "B": 0, "C": 0}
	map3 := map[string]int{"A": 1, "B": 2, "C": 3}

	fmt.Printf("map1 == map2: %t\n", mapComparison(map1, map2)) // false
	fmt.Printf("map1 == map3: %t\n", mapComparison(map1, map3)) // true

}
func mapComparison(map1, map2 map[string]int) bool {
	if len(map1) == len(map2) {
		for key1, value1 := range map1 {
			// If the value for the key of the first doesn't exist in the second one
			if value2, doesExist := map2[key1]; !doesExist || value1 != value2 {
				return false
			}
		}
		// If we've matched all values and got here, return true. Maps have the same pairs
		return true
	}
	return false
}

// INPUT STRINGS DEDUPLICATING
func DeDup() map[string]int {
	seen := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if _, isAlreadyExist := seen[line]; !isAlreadyExist {
			seen[line] += 1
		}
	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "The error occured: %v", err)
		return nil
	}
	return seen
}
func PrintClearData(function func() map[string]int) {
	output := function()

	for _, str := range output {
		fmt.Println(str)
	}
}

// INCOMPARABLE KEYS OF A MAP (SLICES, MAPS AND OTHERS). THE BEST TRICK I'VE EVER SEEN
/*
Add(addCallParameters map[string]int, list []string) increments each string representation of the parameters passed
into Add
*/
func Add(addCallsQuantities map[string]int, list []string) {
	addCallsQuantities[convertListToString(list)]++
}

/*
GetCount(addCallsQuantities map[string]int, list []string) returns the corresponding quantity of Add(...) calls
*/
func GetCount(addCallsQuantities map[string]int, list []string) (int, bool) {
	if callsQuantity, callsExist := addCallsQuantities[convertListToString(list)]; !callsExist {
		return 0, false
	} else {
		return callsQuantity, true
	}
}

/*Convert a slice of string to a string */
func convertListToString(list []string) string {
	return fmt.Sprintf("%q", list)
}

// COUNT OF VARIETY TOKENS
func CharCount(fileName string) {

	quantities := make(map[string]map[rune]int) // Count of Unicode characters
	quantities["Letters"] = make(map[rune]int)
	quantities["Digits"] = make(map[rune]int)
	quantities["Others"] = make(map[rune]int)

	var utflen [utf8.UTFMax + 1]int // Count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	fmt.Println(unicode.Is)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The error: | %v | occured ", err)
		os.Exit(1)
	}
	defer file.Close()

	input := bufio.NewReader(file)
	for {
		/*
			1) Read a single rune, get its UTF-8 representation as []byte, get rune size.
			2) If rune is invalid we get replacement symbol assigned to a run variable
			3) err will be used to recognize the break condition
		*/
		run, amountOfBytes, err := input.ReadRune()

		// If we've reached
		if err == io.EOF {
			// Stop file scanning
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "The error occured: %v\n", err)
			os.Exit(1)
		}

		// If we got an invalid symbol
		if run == unicode.ReplacementChar {
			// Increment the invalid symbols counter
			invalid++
		}

		// Increment a count for each valid unicode symbol
		if unicode.IsLetter(run) {
			quantities["Letters"][run]++
		} else if unicode.IsDigit(run) {
			quantities["Digits"][run]++
		} else {
			quantities["Others"][run]++
		}

		// Increment a count for each bytes length in UTF-8 representation
		utflen[amountOfBytes]++
	}

	// Print quantities for each unicode valid symbol
	for mapKey, subMap := range quantities {
		fmt.Printf("\n------------ %s ------------\n", mapKey)
		for runeKey, runeCount := range subMap {
			fmt.Printf("%5c\t%5d\n", runeKey, runeCount)
		}
	}

	for i := 1; i < len(utflen); i++ {
		fmt.Printf("%d. byte(s) count: %5d\n", i, utflen[i])
	}

	if invalid != 0 {
		fmt.Printf("Replacement symbol count: %d\n", invalid)
	}
}
func WordFrequency(fileName string) {

	wordsQuantities := make(map[string]int) // Count of Unicode characters

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The error: | %v | occured ", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create buffer to store bytes and work with them
	// In this case bytes will be extracted from file
	scanner := bufio.NewScanner(file)

	// Initially we define tokens that will be scanned
	scanner.Split(bufio.ScanWords)

	// Scan the next token(Here's token is represented as word)
	for scanner.Scan() {
		word := scanner.Text() // Get the data of the fetched token
		// Increment its counter in map
		wordsQuantities[word]++
	}

	// Print all the data of the fully filled map
	for key, value := range wordsQuantities {
		fmt.Printf("%10s - %d\n", key, value)
	}
}

// GRAPH
func CreateGraph(countOfPeaks int) map[int]map[int]bool {
	// Create the body of the graph (each peak is mapped to other peaks)
	graph := make(map[int]map[int]bool)
	for i := 1; i <= countOfPeaks; i++ {
		// Create a map for the current peak
		graph[i] = make(map[int]bool)
		for j := 1; j <= countOfPeaks; j++ {
			// Create a connection to another peak randomly
			graph[i][j] = getRandBool()
		}
	}
	return graph
}
func AddEdge(from, to int, graph map[int]map[int]bool) {
	if ok1, ok2 := graph[from] == nil, graph[to] == nil; ok1 {
		peakFrom := make(map[int]bool)
		peakFrom[to] = true
		graph[from] = peakFrom
		if ok2 {
			peakTo := make(map[int]bool)
			graph[to] = peakTo
		}
		return
	} else if !HasEdge(from, to, graph) {
		graph[from][to] = true
	}
}
func HasEdge(from, to int, graph map[int]map[int]bool) bool {
	return graph[from][to]
}
func PrintGraph(graph map[int]map[int]bool) {
	for peak, value := range graph {
		fmt.Printf("peak: %d -> %5v\n", peak, value)
	}
}
func getRandBool() bool {
	return rand.Intn(101) > 50
}
