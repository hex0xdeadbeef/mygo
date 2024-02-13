package chapter7

/*
7. INTERFACES---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Go's interfaces are satisfied implicitly. There's no need to write the all interfaces the type implements externally
	1) There's only the need to possess the necessary methods.

	2) The design let's create new interfaces without changing the existing types adding to them the list of interfaces
	it implements now, which is particularly useful for types defined in packages that we don't control.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.1 INTERFACES AS CONTRACTS---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
1. The interface is an abstract type that only talk about the required behavior of its implementers.

2. An interface reveals only its methods.

3. An intreface defines the contract between itself and its callers.
	1) Contract requires that the caller provide a function with a value of a concrete type that has all the methods an
	interface has.

	2) Contract guarantees that the method has an interface parameter will do its job given any value that
	satisfies the interface of the parameter.

4. We can safely pass a value of any concrete type to the function interface parameter if it satisfies this interface.

5. Declaring a String method makes a type satisfy one of the most widely used interfaces of all: fmt.Stringer

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.2 INTERFACES AS CONTRACTS---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
1. An interface type specifies a set of methods that a concrete type must possess to be considered as an instance of that
interface.

2. There might be some combinations of interfaces in a single interface type. The approach to reach it resembles struct
embedding. It's called "embedding" an interface. The all embedded interface anonymous fields promote their methods requir
ements to the enclosing interface.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.3 INTERFACE SATISFACTION---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
1. A type satisfies an interface if it possesses all the methods the interface requires.

2. The assignability rule: an expression may be assigned to an interface variable only if its type satisfies the interface.
	! This rule applies even when the right-hand side is itself an interface !

3. If the interface is satisfied only for *T receiver there will be no ability to interact with this interface
directly through T instances. We will have to explicitly precede it with "&".

4. Like an envelope that wraps and conceals the letter it holds, an interface wraps and conceals the concrete type descriptor
and value it holds. Only the methods revealed by the interface type may be called, even if the concrete type has others.

5. The type: interface{} is called the "empty interface" type is indispensable. The empty interface type places no
demands on the types that satisfy it. WE CAN ASSIGN/HOLD ANY VALUE TO/IN THE EMPTY INTERFACE VARIABLE. All the types satisfy it.

6. Since interface satisfaction depends only on the methods of the two types involved, there is no need to declare the
relationship between a concrete type and the interfaces it satisfies. That said, it's occasionally useful to document
and assert the relationship.

7. Non-empty interface types are most often satisfied by a pointer type, particularly when one or more of the interface
methods implies some kind of mutation to the receiver. A pointer to a struct is an especially common method-bearing type.

8. A concrete type may satisfy many unrelated interfaces.

9. Each grouping of concrete types based on their shared behaviors can be expressed as an interface type.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.4 PARSING FLAGS WITH flag.Value---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
!CHECK THE CODE!

7.5 INTERFACE VALUES---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Interface value has two components:
	1) A concrete type descriptor
	2) Value of concrete type (point 1) )
It's called "dynamic value"

2. For statically typed languages like Go, types are a compile-time concept, so a type is not a value. The type isn't
contained in memory because it's just a label that tolds what values can be stored.

3. In Go, variables are always initialized to a well-defined value, and interfaces are no exception. The zero value
for an interface has both its type and value components set to nil. Where:
	1) Type is the type descriptor of the implementer
	2) Value is data of the implemetner that is assigned to an interface variable

4. Calling any method of a nil interface value causes a panic.

5. After implicit/explicit conversion from a concrete type operand to an interface variable, the interface variable
captures the type and the value of its operand.

6. Dynamic dispatch:
	1)  The compiler generates code to obtain the address of the method named "MethodName" from the type descriptor,
	then make an indirect call to that address.
	2) The receiver argument for the call is a copy of the interface's dynamic value.

7. An interface value can hold arbitrarily large dynamic values.

8. Interface values can be compared using == and !=. Two interface values are equal if:
	1) Both are nil
	OR
	2) Their dynamic types are identical && their dynamic values are equal according of == for that type.

9. Because interface values are comparable, they may be used as the keys of a map or as the operand of a switch statement.
	! IF TWO INTERFACES VALUES ARE COMPARED AND HAVE THE SAME DYNAMIC TYPE, BUT THAT TYPE ISN'T COMPARABLE (E.G.
	SLICE), THEN THE COMPARISON FAILS WITH A PANIC !

10. When comparing interface values or aggregate types that contain interface values, we must be aware of the potential
for a panic. The same risk exists when using interfaces as map keys or switch operands.
	! WE SHOULD ONLY COMPARE INTERFACES IF WE'RE CERTAIN THAT THEY CONTAIN DYNAMIC VALUES OF COMPARABLE TYPES !

11. When handling errors or during debugging, if's often helpful to report the dynamic type of an interface value.
We can reach it by using "%T" verb.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

7.5.1 CAVEAT: AN INTERFACE CONTAINING A NIL POINTER IS NON-NIL---------------------------------------------------------------------------------------------------------------------------------------------------------------
CHECK THE CODE

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.6 SORTING WITH sort.Interface---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. So that the sequence of type's instances can be sorted by function sort.Sort, the type must possesses three functions:
	1) Len() int
	2) Less (i, j int) bool
	3) Swap(i, j int)

2. Suppose we have []string and want to use sort.Sort. We should make a named type, implement Len(), Less() bool,
Swap() and make explicit conversion TypeName([]string) that yields a slice with the same length/capacity/UNDERLYING
ARRAY []string.

3. Sorting a slice of string is so common that the "sort" package provides the StringSlice type, as well as a function called Strings().

4. To sort in the DESC order we use: sort.Sort(sort.Reverse(data)). sort.Reverse uses the composition

5. We can define our logic with embedding to Less(i, j int) function our own function that determines the order.

6. There is the fuction IsSorted() that checks whether the slice is sorted in the ASC order or not.

7. For convinience, the "sort" package provides versions of its functions and types specialized for:
	1) []int
	2) []string
	3) []float64

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.7 THE http.Handler INTERFACE---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Function ListenAndServe(address string, h Handler) error runs forever or until the server fails (or fails to start) with an error, always
non-nil, which if returns.

2. We can switch in the implementer's ServeHTTP(w http.ResponseWriter, req *http.Request) using switch with req.URL.path as a key. In real
applications it's convinient to define the logic for each case in a separate function or method.

3. Furthermode, related URLs may need similar logic to 2). Several image files may have URLs of the form "/images/*.png", for instance.
For these reasons "net/http" provides ServeMux (request multiplexer), to simplify the association between URLs and handlers.
	1) ServeMux aggregates a collection of http.Handler into a single http.Handler.

	2) So that we can pass a function/method we should explicitly convert it with http.HandleFunc() so that it satisfies the http.Handler
	interface.
		type HandlerFunc func(w ResponseWriter, r *Request)

		func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
			f(w, r)
		}

		! HandlerFunc IS A FUNCTION TYPE THAT HAS METHODS AND SATISFIES AN INTERFACE. THE BEHAVIOR OF ITS ServeHTTP METHOD IS TO CALL THE
		UNDERLYING FUNCTION !

	3) There's the simplification multiplexerName.HandleFunc(patter string, function func(w ResponseWriter, r *Request)) that allows to pass
	the functions/methods directly.

	4) In most programs, one web server is plenty.

	5) It's typical to define HTTP handlers across many files of an application, and it'd be a nuisance if they all had to be registered with
	the application's ServeMux instance. So, for convinience, "net/http" provides a global ServerMux instance called DefaultServeMux and pa
	ckage-level functions called http.Handle and http.HandleFunc, so we needn't pass it to ListenAndServe, nil will do.

	! WEB SERVER INVOKES EACH HANDLER IN A NEW GOROUTINE, SO HANDLERS MUST TAKE PRECAUTIONS SUCH AS locking WHEN ACCESSING VARIABLES THAT
	OTHER GOROUTINES, INCLUDING OTHER REQUESTS TO THE SAME HANDLER, MAY BE ACCESSING !

4. http.ResponseWriter is another interface. it augments io.Writer with methods for sending HTTP response headers.

5. We can write errors codes to the clients using w.WriteHeader(http.Status...). This might be done before writing any data to w.

6. We can pass errors messages to a client using fmt.Fprintf(responseWriter, format, ...).

7. There's the equal to 5+6 steps function: Error(w ResponseWriter, error string, code int). It sets the header, writes a given message to the client side.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.8 THE error INTERFACE---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. The error is just an interface type with a single method that returns an error message.
	type error interface {
		Error() string
	}

2. The simplest way to create an error is by calling errors.New(), which returns a new error for a given error message.

3.
package error

func New(text string) error { return &errorString{text} }

type errorString struct {
	text string
}

func (e *errorString ) Error() string {
	return e.text
}

	1) The underlying type of errorString is a struct, not a string, to protect its representation from inadvertent(or premediated) updates
	2) The reason of the pointer type *errorString satisfies the error interface so that every call to New allocates a distinct error instance
	that is equal to no other.

4. Calls to errors.New is relatively infrequent because there's a convenient wrapper function: fmt.Errorf() that does string formatting too.
	package fmt

	import "errors"

	func Errorf(format string, args ... interface{}) error {
			return errors.New(Sprintf(format, args...))
	}

5. "syscall" package provides Go's low-level system call API. It provides programmers with system errors as well. It defines a numeric type
"Errno" that satisfies error interface, and on Unix platforms, Errno's Error() method does a lookup in a table of strings.

6. Errno is an efficient representation of system call errors drawn from a finit set, and it satisfies the standart error interface.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.9 THE error INTERFACE---------------------------------------------------------------------------------------------------------------------------------------------------------------
CHECK THE CODE
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.9 TYPE ASSERTIONS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. An assertion is an operation applied to an interface value. Syntactically it looks like x.(T) where "x" is an expression of
an interface type and "T" is a type called the "asserted" type.

2. A type assertion checks that the dynamic type of its operands matches the asserted type. There are two possibilities:
	1) If the asserted type "T" is a concrete type, then the type assertion checks whether x's dynamic type is identical to T.
		1) If this check succeeds, the result of the type assertion is x's dynamic value, whose type is of course T.
		2) If this check fails, then the operation panics.

	2) If instead the asserted type T is an interface type, then the type assertion checks wheter x's dynamic type satisfies T.
		1) If this check succeeds, the dynamic value isn't extracted. The result is still an interface value with the same type and
		value components, but the result has the interface type T.
		In other words, a type assertion to an interface type changes the type of the expression, making a different (and usually larger)
		set of methods accessible, but it preserves the dynamic type and value components inside the interface value.

		2) No matter what type was asserted, if the operand is a nil interface value, the type assertions fails.

		3) A type assertion to a less restrictive interface type (one with fewer methods) is rarely needed, as it behaves just like
		an assignment, except the nil case ( 2) )

3. Often we're not sure of the dynamic type of an interface value, and we'd like to test whether it's some particular type. So that
check the validity of a type we use the following logic:

	var w io.Writer = os.Stdout
	f, ok := w.(*os.File) // success: ok, f == os.Stdout
	b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil

if the test is passed, "ok" is given true and the left side variable is accepted the dynamic value and the type of an interface variable, otherwise
the "ok" variable takes false and the left side variable turns into nil value.
	1) The second result is conventionally assigned to a variable named "ok"

	2) The result is often immediately used to decide what to do next. The extended form of the if statement makes this quiet compact.

	3) When the operand of a type assertion is a variable, rather than invent another name for the new local variable, you'll sometimes
	see the original name reused, shadowing the original.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.11 DISCRIMINATING ERROR WITH TYPE ASSERTIONS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. The "os" package defines a type called "PathError" to describe failures involving an operation on a file path, like open or delete
and a variant called LinkError to desctibe failures of operations involving two file paths.
	package os
	type PathError struct {
		Op string
		Path string
		Err error
	}

	func (e *PathError) Error() string {
		return e.Op + " " + e.Path + ": " + e.Err.Error()
	}

2. Clients that need to distinguish one kind of failure from another can use a type assertion to detect the specific of the error. The
specific type provides more detail than a simple string.

3. Error discrimination must usually be done immediately after the failing operation, before an error is propagated to the caller.
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

7.12 QUERYING BEHAVIORS WITH INTERFACE TYPE ASSERTIONS---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Consider we have a consumer that has a crate which can be filled with a data using method of the interface io.Writer or any else,
but the parameter of enclosing function is io.Writer that only has the method Write:

	Wtite(p []byte) (n int, err error)

but we want put the strings into the crate because conversion from a string to the corresponding byte slice isn't efficient because of
allocating new data and copying it. And the argument has the method WriteString(s string). What should we do in order to avoid the tem
porary copying?

2. We cannot assume that an arbitrary io.Writer "w" also has the WriteString() method, but we can create the new interface that has only
the method WriteString() and use assertion to test whether the dynamic type of "w" satisfies this new interface.

3. The technique above relies on the assumption that IF a type satisfies the interface below, THEN WriteString(s) must have the same
effect as Write([]byte(s))
	interface {
		io.Write
		WriteString(s string) (n int, err error)
	}

4. Assumption like the above one should be documented properly, so that clients won't be confused.

5. The WriteString() function above uses a type assertion to see whether a value of a general interface type also satisfies a more
specific interface type, and if so, it uses the behaviors of the specific interface. It's the same how fmt.Fpringf distinguishes va
lues that satisfy error or fmt.Stringer from all other values. Within fmt.Fprintf, there's a step that converts a single operand to
a string, something like this:

	package fmt

	func formatOneValue(x interface{}) string {
		if err, ok := x.(Stringer); ok {
			return str.String()
		}
		...
	}

5. To avoid repeating ourselves we can move the check into the utility function.

6. The standart library provides the function io.WriteString() and it's recommended way to write a string to an io.Writer.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.13 TYPE SWITCHES---------------------------------------------------------------------------------------------------------------------------------------------------------------
1. Interfaces are used in two distinct styles:

	SUBTYPE POLYMORPHISM
	1) THE METHODS EMPHASIS. Interfaces express the similarities of the concrete types that satisfy the interface but hide
	the representation details and intrinsic operations of those concrete types.

	AD HOC POLYMORPHISM
	2) THE ARBITRARY TYPES HOLDING EMPHASIS. The second style exploits the ability of an interface values to hold values of a variety
	of concrete type and considers the interface to be the union of those types. Type assertions are used to discriminate among these
	types dynamically and treat each case differently. In this case, the emphasis is on the concrete types that satisfy the interface,
	not on the interface's methods (if indeed it has any), and there's no hiding of information. These interfaces are called as "disc-
	riminated unions"

2. Instead of the consequent checks in "if" statements we can switch of the values of interface{} function parameter. It simplifies
an if-else chain that performs a series of value equality tests. An analogous "type switch" statement simplifies an if-else chain of
type assertions.
	1) x.(type) - that's literally the keyword type
	2) Each case has one or more types.
	3) A type switch enables a multi-way branch based on the interface value's dybamic type.
	4) The nil case matches if x == nil
	5) default case matches if no other case does.

3. As with an ordinary switch statement, cases are considered in order and, when a match is found, the case's body is executed.

4. Case order becomes significant when one or more case types are interfaces, since the there's a possibility of two cases matching.

5. The position of the "default" case relative to the others is immaterial.

6. No fallthrough is allowed.

7. We can extend the switch statement in the situation when we need the value of a parameter to be processed. The extension is:
	switch x := x.(type) {
		...
	}

	1) In this case we reuse "x" name to use it in switch statement. Because a switch implicitly creates a new lexical block the decla
	ration of the new variable called x doesn't conflict with a variable x in an outer block.

	2) Each case also implicitly creates a separate lexical block.

	3) In this version, within the block of each single-type case, the variable x has the same type as the case.

8. Although the type of x is interface{}, we consider it a discriminated union of the types are placed in cases.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


7.14 EXAMPLE: TOKEN-BASED XML DECODING---------------------------------------------------------------------------------------------------------------------------------------------------------------
XML - Extensible Markup Language

The example of XML data:
<bookstore> // <Name>
    <book category="cooking"> <Name Attr.Name = Attr.Value>
        <title lang="en">Everyday Italian</title>
        <author>Giada De Laurentiis</author>
        <year>2005</year>
        <price>30.00</price>
    </book>
    <book category="children">
        <title lang="en">Harry Potter</title>
        <author>J.K. Rowling</author>
        <year>2005</year>
        <price>29.99</price>
    </book>
</bookstore>

package xml

type Name struct {
	Local string // e.g. "Title", "id"
}

type Attr struct { // e.g. name="value"
	Name Name
	Value string
}

type Token interface{} // includes StartElement, EndElement, CharData and Comment

type StartElement struct { // e.g. <name>
	Name Name
	Attr []Attr
}

type EndElement struct { // e.g. </name>
	Name Name
}

type CharData []byte // <p>CharData</p>

type Comment []byte // e.g. <!-- Comment -->

type Decoder struct { ... }

func NewDecoder(io.Reader) *Decoder
func (*Decoder) Token() (Token, error) // returns next Token in sequence

1. The "encoding/xml" package proved a similar to the "encoding/json" API. This API is convenient when we want to construct a
representation of the document tree, but that's unnecessary for many programs.

2. The "encoding/xml" package also provides a lower-level token-based API for decoding XML. In the token-based style, the parser con
sumes the input and produces a stream of tokens, primarily of four kinds: "StartElement", "EndElement", "CharData", "Comment".

3. Each call to (*xml.Decoder).Token returns a token.

4. The "Token" interface, which has no methods, is also an example of a descriminated union.
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

*/
