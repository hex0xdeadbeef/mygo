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

2. The function "Sprint()" supports the narrow range of cases that return a string representation of the value. But how do we deal with other types, like []float64, map[string][]string and
so on? And what about named types, like url.Values? Without a way to inspect the representation of values of unknown types, we quickly get stuck. What we need is reflection.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
12.2 reflect.Type AND reflect.Value

1. Reflection is provided by "reflect" package. It defines two important types:
	1) "Type" represents a Go type. It's an interface with many methods for discriminating among types and inspecting their components, like the fields of a struct or the parameters of a
	function. The sole implementation of "reflect.Type" is the type descriptor, the same entity that identifies the dynamic type of an interface value.
	2) "Value" can hold a value of any type.

2. The "reflect.TypeOf()" function accepts any "interface{}" values and returns its dynamic type as a "reflect.Type". Because "reflect.TypeOf()" returns an interface value's dynamic type,
it always returns a concrete type, but it's capable of referencing interfaca types too.
	1) "reflect.Type" satisfies "fmt.Stringer". Because printing the dynamic type of an interface value is useful for debugging and logging, "fmt.Printf()" provides the shorthand "%T" that
	uses "reflect.TypeOf() internally.

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
	1) Slices and arrays: the logic is the same for both. The "Len()" method returns the number of elements of a slice/array and "Index(i)" retrieves the element at index i, also reflect.Value.
	"Index(i)" panics if "i" is out of bounds. These are analogous to the built-in len(a) and a[i] operations on sequences. The display function recursively invokes itself each element of the
	sequence, appending the subscript notation "[i]" to the path.

	Although "reflect.Value" has many methods, only a few are safe to call on any given value. For example, the Index method may be called on values of kind Slice/Array/String, but panics for
	any other kind.

	2) Structs: The "NumField()" method reports the number of fields in the struct, and "Field(i)" returns the value of the i-th field as a "reflect.Value". The list of fields includes ones
	promoted from anonymous fields. To append the field selector notation ".f" to the path, we must obtain the "reflect.Type" of the struct and access the name of i-th field.

	3) Maps: The "MapKeys()" returns a slice of "reflect.Value" values one per map key. As usual when iterating over a map, the order is undefined. "MapIndex(key)" returns the value
	corresponding
	to key. We append the subscript notation "[key]" to the path. (We're cutting a corner here. The type of a map key isn't restricted to the types formatAtom handles best; arrays, structs, and
	interfaces can also be valid map keys.)

	4) Pointers: The "Elem()" method returns the variable pointed to by a pointer, again as a "reflect.Value" This operation would be safe even if the pointer is nil, in which case the result
	would have kind "Invalid", but we use "IsNil()" to detect nil pointers explicitly so we can print a nore appropriate message. We prefix the path with a "*" and parenthesize it to avoid
	ambiguity.

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
	1) In the first example, Display calls "reflects.ValueOf(number)", which returns a value of kind "Int". As we mentioned "reflect.ValueOf()" always returns a value of a concrete type since
	it extracts the contents of an interface value.
		Display number (int):
		number = 3

	2) In the second example, "Display()" calls "reflect.ValueOf(&number)", which returns a pointer to "number", of kind Ptr. The switch case for Ptr calls Elem on this value, which returns a
	Value representing the variable "number" itself, of kind Interface. A Value obtained indirectly, like this one, may represent any value at all, including interfaces.
		Display &number (*interface {}):
		(*&number).type = int
		(*&number).value = 3

4. As currently implemented, "Display()" will never terminate if it encounters a cycle in the object graph, such as linked list that eats its own tail.
	1) Cycles pose less of a problem for "fmt.Sprint()" because it rarely tries to print the complete structure. For example, when it encounters a pointer, it breaks the recursion by printing
	the pointer's numeric value.
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
	IN FACT NO "reflect.Value" RETURNED BY "reflect.ValueOf(x)"

	4) In contrast, "d" derived from by dereferencing the pointer within it, refers to a variable and is thus addressable. The operation "Elem()" is the same to dereferencing. We obtain
	an addressable an addressable Value for any variable x.

4. We can ask a "reflect.Value" whether it's addressable through its "CanAddr()" method.

5. To recocer the variable from an addressable "reflect.Value" requires three steps:
	1) We call "Addr()" which returns a "Value" holding a pointer to the variable.
	2) We call "interface()" on this value, which returns an interface{} value containing the pointer.
	3) If we know the type of the variable, we can use a type assertion to retrieve the contents of the interface as an ordinary pointer.

6. We can update the variable referred to by an addressable "reflect.Value" directly, without using a pointer, by calling the "reflect.Value.Set" passing the argument using "reflect.
wValueOf()".
	1) The program will panic if the types of the underlying type and the passed argument's type don't matches
	2) The program panic if we try to use the "Set()" method on a unaddressableValue variable as well.

There are variants of "Set" specialized for certain groups of basic types, they look like "SetXxx()".
	1) In some ways these methods are more forgiving. SetInt for example, will succeed so long as the variable's type is some kind of integer, or even a named type whose underlying type
	is a signed integer. If even a value is too large it'll be quietly truncated to fit.
	2) Calling "SetInt()" or just "Set()"  on a "reflect.Value" that refers to an "interface{}" variable will panic.
	3) If an addressable variable points to "interface{}" value we can pass any values to "Set(reflect.Valueof(...))", but using the type specified methods "Value.SetXxx(...)" on an
	addressable variable will cause panic, because they're supposed to be used on the particular type.

7. Reflection can read unexported components / inspect their values but it cannot change them.
	1) So that we can find out whether the field is settable or not, we should use "reflect.Value.CanSet(...)" function, that returns a corresponding boolean value.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
