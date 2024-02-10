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

	2) So that we can pass the function/method we should explicitly convert it with http.HandleFunc() so that it satisfies the http.Handler
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

7. There's the equal to 5+6 steps function:Error(w ResponseWriter, error string, code int). It sets the header, writes a given message to the client side.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
*/
