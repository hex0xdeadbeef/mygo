package chapter12

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12. REFLECTION

1. Reflection is needed because it allows:
	1) Update variables and inspect their values at run time.
	2) Call variables' methods
	3) Apply the operations intrinsic to variables' representation.
	4) Treat types themselves as first-class values.

	ALL WITHOUT KNOWING THEIR TYPES AT COMPILE TIME.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.1 WHY REFLECTION?

1. Sometimes we need to write a function capable of dealing uniformly with values of types that:
	1) Don't satisfy a common interface
	OR
	2) Don't have a known representation
	OR
	3) Don't exist at the time we design the function
	OR EVEN 1), 2), 3)

A familiar example is the formatting logic within "fmt.Fprintf()", which can usefully print an arbitrary value of any type, even a user defined one.

2. The function "Sprint()" supports the narrow range of cases that return a string representation of the value. But how do we deal with other types, like []float64, map[string][]string and so on? And what about named types, like "url.Values"? Without a way to inspect the representation of values of unknown types, we quickly get stuck. What we need is reflection.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.2 reflect.Type AND reflect.Value

1. Reflection is provided by "reflect" package. It defines two important types:
	1) "Type" represents a Go type. It's an interface with many methods for discriminating among types and inspecting their components, like the fields of a struct or the parameters of a function. The sole implementation of "reflect.Type" is the type descriptor, the same entity that identifies the dynamic type of an interface value.
	2) "Value" can hold a value of any type.

2. The "reflect.TypeOf()" function accepts any "interface{}" values and returns its dynamic type as a "reflect.Type". Because "reflect.TypeOf()" returns an interface value's dynamic type, it always returns a concrete type, but it's capable of referencing interfaca types too.
	1) "reflect.Type" satisfies "fmt.Stringer". Because printing the dynamic type of an interface value is useful for debugging and logging, "fmt.Printf()" provides the shorthand "%T" that uses "reflect.TypeOf() internally.

3. The "reflect.ValueOf()" function accepts any "interface{}" and returns a "reflect.Value" containing the interface's dynamic value. As with "reflect.TypeOf()", the results of
"reflect.ValueOf()" are always concrete, but a "reflect.Value" can hold interface values too.
	1) "reflect.Value" also satisfies fmt.Stringer(), but unless Value holds a string, the result of the String method reveals only the type. Instead, use the fmt package's "%v" verb
	which treats "reflect.Value" specially.
	2) Calling the "Type()" method on a "Value" returns its type as a "reflect.Type"
	3) There's inverse operation to "reflect.ValueOf()" - "reflect.Value.Interface()". It returns an "interface{}" holding the same concrete value as the "reflect.Value"

4. A "reflect.Value" and an "interface{}" can both hold arbitrary values. Their differencies:
	1) An empty interface hides the representation and intrinsic operations of the value it holds and exposes none of this methods, so unless we know its dynamic type and use a type
	assertion to peer inside it, there's a little we can do to the value within.
	2) Value has many methods for inspecting its contents, regardless of its type.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.3 Display, RECURSIVE VALUE PRINTER
1.
	1) SLICES and ARRAYS: the logic is the same for both. The "Len()" method returns the number of elements of a slice/array and "Index(i)" retrieves the element at index i, also
	"reflect.Value". "Index(i)" panics if "i" is out of bounds. These are analogous to the built-in len(a) and a[i] operations on sequences. The display function recursively invokes
	itself each element of the sequence, appending the subscript notation "[i]" to the path.

	Although "reflect.Value" has many methods, only a few are safe to call on any given value. For example, the "Index()" method may be called only on values of kind Slice/Array/
	String, but panics for any other kind.

	2) STRUCTS: The "NumField()" method reports the number of fields in the struct, and "Field(i)" returns the value of the i-th field as a "reflect.Value". The list of fields
	includes ones promoted from anonymous fields. To append the field selector notation ".f" to the path, we must obtain the "reflect.Type" of the struct and access the name of i-th
	field.

	3) MAPS: The "MapKeys()" returns a slice of "reflect.Value" values one per map key. As usual when iterating over a map, the order is undefined. "MapIndex(key)" returns the value
	corresponding to key. We append the subscript notation "[key]" to the path. (We're cutting a corner here. The type of a map key isn't restricted to the types "formatAtom()" handles
	best; arrays, structs, and interfaces can also be valid map keys.)

	4) Pointers: The "Elem()" method returns the variable pointed to by a pointer, again as a "reflect.Value" This operation would be safe even if the pointer isn't nil, in which case
	the result would have kind "reflect.Invalid", but we use "IsNil()" to detect nil pointers explicitly so we can print a more appropriate message. We prefix the path with a "*" and
	parentheses it to avoid ambiguity.

	5) Interfaces: again, we use "IsNil()" to test whether the interface is nil, and if not, we retrieve its dynamic value using v.Elem() and print its type and value.

2. We can use "Display()" to display the internals of library types such as *os.File:
	(*(*strangelove).file).pfd.fdmu.state = 0
	(*(*strangelove).file).pfd.fdmu.rsema = 0
	(*(*strangelove).file).pfd.fdmu.wsema = 0
	(*(*strangelove).file).pfd.Sysfd = 2
	(*(*strangelove).file).pfd.SysFile.iovecs = nil
	(*(*strangelove).file).pfd.pd.runtimeCtx = 0
	(*(*strangelove).file).pfd.csema = 0
	(*(*strangelove).file).pfd.isBlocking = 1
	(*(*strangelove).file).pfd.IsStream = true
	(*(*strangelove).file).pfd.ZeroReadIsEOF = true
	(*(*strangelove).file).pfd.isFile = true
	(*(*strangelove).file).name = "/dev/stderr"
	(*(*strangelove).file).dirinfo = nil
	(*(*strangelove).file).nonblock = false
	(*(*strangelove).file).stdoutOrErr = true
	(*(*strangelove).file).appendMode = false

NOTICE THAT EVEN UNEXPORTED FIELDS ARE VISIBLE TO REFLECTION.

3. "DifferencePtrAndNonPtrDisplaying()":
	1) In the first example, Display calls "reflects.ValueOf(number)", which returns a value of kind "reflect.Int". As we mentioned "reflect.ValueOf()" always returns a value of a
	concrete type since it extracts the contents of an interface value.
		Display number (int):
		number = 3

	2) In the second example, "Display()" calls "reflect.ValueOf(&number)", which returns a pointer to "number", of kind "Ptr". The switch case for "reflect.Ptr" calls Elem on this
	value, which returns a Value representing the variable "number" itself, of kind Interface. A Value obtained indirectly, like this one, may represent any value at all, including
	interfaces.
		Display &number (*interface {}):
		(*&number).type = int
		(*&number).value = 3

4. As currently implemented, "Display()" will never terminate if it encounters a cycle in the object graph, such as linked list that eats its own tail.
	1) Cycles pose less of a problem for "fmt.Sprint()" because it rarely tries to print the complete structure. For example, when it encounters a pointer, it breaks the recursion by printing the pointer's numeric value.
	2) It can get stuck trying to print a slice or map that contains itself as an element, but such rare cases don't warrant the considerable extra trouble of handling cycles.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.5 SETTING VARIABLES with reflect.Value
1. x, x.f[1], *p denote variables, but x + 1, f(2) don't.

2. All the usual rules for addressability have analogs for reflection.

3. A variable is an addressable storage location that contains a value, and its value may be updated through that address. A similar distinction applies to "reflect.Value"s. Some
are addressable, others are not.

	x := 2						value	type	variable?
	a := reflect.ValueOf(2)		2		int		no
	b := reflect.ValueOf(x)		2		int		no
	c := reflect.ValueOf(&x)	&x		*int	no
	d := c.Elem()				2		int		yes (x)

	1) The value "a" isn't addressable, it's merely copy of 2
	2) The same is true for "b"
	3) The value within "c" is also non-addressable, being a copy of the pointer value "&x"
	IN FACT NO "reflect.Value" RETURNED BY "reflect.ValueOf(x)" IS ADDRESSABLE

	4) In contrast, "d" derived from by dereferencing the pointer within it, refers to a variable and is thus addressable. The operation "Elem()" is the same to dereferencing. We obtain an addressable Value for any variable x.

4. We can ask a "reflect.Value" whether it's addressable through its "CanAddr()" method.

5. To recover the variable from an addressable "reflect.Value" requires three steps:
	1) We call "Addr()" which returns a "Value" holding a pointer to the variable.
	2) We call "interface()" on this value, which returns an interface{} value containing the pointer.
	3) If we know the type of the variable, we can use a type assertion to retrieve the contents of the interface as an ordinary pointer.

6. We can update the variable referred to by an addressable "reflect.Value" directly, without using a pointer, by calling the "reflect.Value.Set(...)" passing the argument using
"reflect.ValueOf()".
	1) The program will panic if the types of the underlying type and the passed argument's type don't match.
	2) The program panic if we try to use the "Set()" method on an unaddressable "Value" variable as well.

There are variants of "Set" specialized for certain groups of basic types, they look like "SetXxx()".
	1) In some ways these methods are more forgiving. SetInt for example:
		1)  Will succeed so long as the variable's type is some kind of integer
		OR
		2) Even a named type whose underlying type is a signed integer.
		OR
		3) If even a value is too large it'll be quietly truncated to fit.
	2) Calling "SetInt()" or just "Set()" on a "reflect.Value" that refers to an "interface{}" variable will panic.
	3) If an addressable variable points to "interface{}" value we can pass any values to "Set(reflect.Valueof(...))", but using the type specified methods "Value.SetXxx(...)" on an
	addressable variable will cause panic, because they're supposed to be used on the particular type.

7. Reflection can read unexported components / inspect their values but it cannot change them.
	1) So that we can find out whether the field is settable or not, we should use "reflect.Value.CanSet(...)" function, that returns a corresponding boolean value.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.6 DECODING S-EXPRESSION
1. The "Unmarshall()" function uses reflection to modify the fields of the existing variable, creating new maps/structs/slices as determined by the corresponding type and the content of the
incoming data.

2. The lexer uses the "Scanner" type from "text/scanner" package to break an input into a sequence of tokens such as comments/identifiers/string literals/numeric literals. The scanner's
"Scan()" method advances the scanner and returns the kind of the next token, which has type "rune", but the "text/scanner" package represents the kinds of multicharacters tokens "Ident"/
"String"/"Int" using small negative values of type "rune". Following a call to "Scan()" that returns one of these kinds of token, the scanner's "TokenText" method returns the text of the
token.
	1) Since a typical parser may need to inspect the current token several times, but the "Scan()" method advances the scanner, we wrap the scanner in a helper type, called "lexer" that
	keeps track of the token most recently returned by "Scan()".

3. Our S-expressions use identifiers for two distinct purposes:
	1) Struct field names
	2) nil value for a pointer. The "read(...)" function handles only this case. When it encounters the "scanner.Ident" "nil", it sets "v" to zero value of its type. For any other identifier,
	it reports an error.

The readList function handles identifiers used as struct field names.

4. A '(' token indicates the start of a list.

5. The second function, "readList()", decodes a list into a variable of composite type - a struct/map/slice/array - depending on what kind of Go variable we're currently populating. In
each case, the loop keeps parsing items until it encounters the matching close parenthesis, ')', as detected by "the endList()" function.

6. The interesting part is the recursion. The simplest case is an array.
	1) Until the closing ')' is seen we use "Index()" to obtain the variable for each array element and make a recursive call to "read()" to populate it. As in many error cases, if the
	input data causes the decoder to index beyond the end of array, the decoder panics.
	2) The similar approach is used for slices, except we must create a new variable for each element, populate it, the append it to slice.
	3) The loops for structs/maps must parse a (key value) sublist on each iteration.
		1) For structs, the key is a symbol identifying the field. Analogous to the case for arrays, we obtain the existing struct field using "FieldByName()" and make a recursive call to
		populate it.
		2) For maps, the key may be of any type, and analogous to the case of slices, we create a new variable, recursively populate it, and finally insert the new key/value pair into a
		map.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.7 ACCESSING STRUCT FIELD TAGS
1. In a web server, the first thing most HTTP handler functions do is extract the request parameters into local variables.

2. The "Unpack()" function below does three things:
	1) It calls "req.ParseForm()" to parse the request.
	2) Thereafter, "req.Form" contains all the parameters, regardless of whether the HTTP client used the "GET" or the "POST" request method.
	3) "Unpack()" builds a mapping from the effective name of each field to the variable for that field. The effective name may differ from the actual name if the field has a tag.
		1) The "reflect.Type.Field()" method returns a "reflect.StructField" that provides information about the type of each field such as its: name, type, optional tag.
		2) The "Tag" field is a "reflect.StructTag", which is a string type that provides "Get()" method to parse and extract the substring for a particular key, such as: http:"..." in this
		case.
	4) Finally "Unpack()" iterates over the name/value pairs of the HTTP parameters and updates the corresponding struct fields. Recall that the same parameter name may appear more than
	once. if this happens, and the field is a slice, then all the values of that are accumulated into the slice. Otherwise, the field is repeatedly overwritten so that only the last value has any affect.
	5) The "populate" function takes care of setting a single field v (or a single element of a slice field) from a parameter value.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.8 DISPLAYING THE METHODS OF A TYPE
1. Both "reflect.Type" and "reflect.Value" have a method called "Method()".
	1) Each "t.Method(i)" call returns an instance of "reflect.Method", a struct type that describes the name and type of a single method
	2) Each "v.Method(i)" call returns a "reflect.Value" representing a method value, that is, a method bound to its receiver. Using the "reflect.Value.Call()" method, it's possible to call
	"Value"s of kind "reflect.Func"
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.9 A WORD OF CAUTION
THREE REASONS TO USE REFLECTION WITH CARE:
1. Reflection-based code can be fragile. For every mistake that would cause a compiler report a type error, there's a corresponding way to misuse reflection, but whereas the compiler reports the
mistake at build time, a reflection error is reported during execution as a panic, possibly long after the program was written or even long after it has started running.

We should keep track of the type, addressability and settability of each reflect.Value. The best way to avoid this fragility is to ensure that the use of reflection is fully encapsulated within
our package and, if possible, avoid "reflect.Value" in favor of specific types in our package's API, to restrict inputs to legal values. If this isn't possible, perform additional dynamic checks
before each risky operation. As an example from the standart library, when "fmt.Printf" applies a verb to an inappropriate operand, it doesn't panic mysteriously but prints an informative error
message.

Reflection also reduces the safety and accuracy of automated refactoring and analysis tools, because they can't determine or rely on type information.

2. Since types serve as a form of documentation and the operations of reflection cannot be subject to static type checking, heavily reflective code is often hard to understand. We should document
the expected types and other invariants of functions that accept an "interface{}" an "reflect.Value" values.

3. Reflection-based function may be one or two orders of magnitude slower than code specialized for a particular type. It's fine to use reflection when it makes the program clearer. Testing is a
particularly good fit for reflection since most test use small data sets, but for functions on the critical path, reflection is best avoided.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
