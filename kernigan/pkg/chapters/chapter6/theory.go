package chapter6

/*
6.1 METHOD DECLARATIONS--------------------------------------------------------------------------------------------------
1. The parameter before function name attaches the function to the type of that parameter. This approach is subtle
because each call to the method results in a creation of parameters copies. func (Type) funcName(...) {}

2. Me can use the extra left parameter in a function body. thus we refer to the caller object.
	1) Extra parameter is called the method's receiver

3. In Go we don't use "this"/"self" keys for the receiver. We just choose receiver names just we would for any
other parameters.
	1) It's good to choose a short name so that it'll be used frequently and be consistent across methods.
	2) The common way to do it is just take a first letter of the type.

	4. There will be no name conflict if in a source file a method of type and a particular function with the same
names are declared. Because of the difference of the names:
	1) Package level function: Distance()
	2) Object method: Point.Distance()

5. The expressions p.Distance/p.X/p.Y are called a selectors, because they select the appropriate method/field
for the receiver p of type Point.

6. Since there is a field with the name X the declaration the method with the same name will cause an ambiguity.
Compiler detects and rejects this situation.

7. Since each type has its own name space for method, we can use the name Distance for other types so long, as
they belong to different types.

8. WE CAN CREATE NAMED SLICE TYPES AS type Path []Point AND BIND TO THEM SOME METHODS

9. Methods may be declared on any type defined in the same package so long as its type is neither an interface nor
a pointer.

10. The compiler determines which function to call based on the method name and the type of the receiver.

11. All the methods are bound to the particular type must have different names, but different types can use the
same name for a method.

12. Methods benefits:
	1) Method names can be shorter. This benefit is magnified for calls ordinating outside the package, since
	they can use the shorter name and omit the package name.

13. Method may be called by a type literal.

6.2 METHODS WITH A POINTER RECEIVER----------------------------------------------------------------------------------------------
1. Instead of creation the copies of method parameters, we can declare an extra left parameter as a pointer of its type.
This approach let us save the memory and change the caller object fields. (*Type) func funcName(...) ... {}

2. The rule is: IF ANY METHOD OF Point HAS A POINTER RECEIVER, THE ALL METHODS OF Point SHOULD HAVE A POINTER RECE
IVER, EVEN ONES DON'T STRICTLY NEED IT.

3. Method declarations aren't permitted on named types, that are themselves pointer/interface types.

4. In order to use methods with an extra pointer receiver we should call this methods for an pointer object of this type. There are
three options to do that:
	EXPLICIT

	1) Work with an object with type *Point
					- - - -
					| 2) Create a pointer to an object with type Pointer
	Ungainly cases:	|
					| 3) Make place-in pointer to an object: (&p).ScaleBy(...)
					- - - -

	IMPLICIT
	4) Use non-pointer object and it'll be implicitly turned into the pointer by the compiler. E.g.
		p := {1,2}
		p.ScaleBy(...) // Here's the implicit conversion to a pointer while calling the method.
	It works only for named objects (variables) including struct fields like p.X.
	FOR LITERALS IT DOESN'T WORK BECAUSE THERE IS NO WAY TO OBTAIN THE ADDRESS OF A TEMPORARY VALUE.

5. Pointer objects can call type methods that have non-pointer receiver 'cause of the ability to take the value.
The compiler just inserst an implicit "*" for us and the value is passed to a method.

4|5 Sum. Either the receiver has the same type as the receiver parameter

		 Or the receiver argument is a variable of type T and receiver parameter has type *T. The compiler
		 implicitly takes the address of the variable

		 Or the receiver argument has type *T and the receiver parameter has type T. The compiler implicitly derefe
		 rences the receiver, in other words, loads the value.

6. If all the methods of a named type T have a receiver type of T itself, it's safe to copy instances of that type.
Calling any of its methods necessarily makes a copy.

7. If any method has a pointer receiver we should avoid copying instances of T because doing so may violate internal
invariants.

6.2.1 nil IS A VALID RECEIVER VALUE----------------------------------------------------------------------------------------------
1. Just as some functions can allow nil pointers as arguments, so do some methods for their receiver, especially if
nil is a meaningful zero value of its type.

2. If we define a type whose methods allow nil as a receiver value, it's worth pointing this out explicitly in its
documentation comment.

3. IF A FUNCTION/MEHOD HAS A POINTER PARAMETER AND WE CALL IT WITH AN ARGUMENT AND SUBSEQUENTLY CHANGE THE VALUE WITH THIS
POINTER AND FINALLY ASSIGN TO THE LOCAL POINTER NIL/OTHER VALUE, WE JUST CHANGE THE DIRECTION OF THE LOCAL POINTER, NOT THE INIT
IAL ONE.

6.3 COMPOSING TYPES BY STRUCT EMBEDDING--------------------------------------------------------------------------------------------------------------------
1. We can select as methods as fields and vice versa if 'course they're definition allows to do it (fields/methods
have a capitalized first letter). Due to this fact we can call methods of the embedded anonymous field Point using
a receiver of type ColoredPoint.

2. The embedded anonymous fields promote theirs accessible methods to struct type they're embedded in.

3. Due to the second point we can construct complex types with COMPOSITION of types.

4. There's no parents and children. So we should make wrappers of the embedded anonymous fields methods explicitly.

5. The type of an anonymous field can be a pointer to a named type, in which case fields and methods are promoted in
directly from the pointed-to object.

6. A struct type may have more than one anonymous field, so this struct have all the embedded fields methods.

7. When a compiler resolves a selector to a method, it first looks for a directly declared method by its name, then for
methods promoted from embedded anonymous fields and so on.
! IF THERE'S A METHODS AMBIGUITY ON THE PARTICULAR LEVEL COMPILER INFORMS ABOUT IT !

8. Methods can be declared only on named types and pointers to them, but thanks to embedding, it's possible and sometimes
useful for unnamed struct types to have methods too.

6.4 METHOD VALUES AND EXPRESSIONS--------------------------------------------------------------------------------------------------------------------
1. Usually we select and call a method in the same expression, but we can separate these steps. Selector p.Distance yields
a method value, a function that binds a method to a specific receiver value p. This function can then be invoked without a
receiver value. In other words we can extract a method and assign it to a variable and subsequetly call when we need it.

2. If we extract the method value and the method's receiver isn't a pointer, any object changes won't be reflected on the
result of the extracted method call. When we extract the method value, it saves a snapshot of the object, because the current
value of object is passed into it.

3. Method values are useful when a package's API calls for a function value, and the client's desired behaviour for that fu
nction is to call a method on a specific receiver.

4. We can also get a method expression, but in this case extracting it we shoud provide the statement with the receiver of
method's receiver type. In fact it looks like: "T.f" or "(*T).f. This expression yields a function value with a regular first
parameter taking the place of the receiver, so it can be called in the usual way.

5. Methods expressions are useful when we need a value to represent a choice among several methods belonging to the same type
so that we can call the chosen one with many different receivers.

6.6 ENCAPSULATION--------------------------------------------------------------------------------------------------------------------
1. If variable/method is inaccessible to a client of the object, it's said that it's encapsulated.

2. Encapsulation is sometimes called "information hiding".

3. Capitalized identifiers (of objects/methods) are exported from the package it's defined in and otherwise they'
re not. This mechanism also limits the access to struct's fields and type's methods.

4. Encapsulation allows us to forbid a direct access to process an object data and do it only by using the API.

5. In Go encapsulation mechanism is different from other languages because encapsulation constraint is bound to
the package. Put another way, encapsulation defines accessibility from other packages.

6. The fields of a struct type are visible to all code within the same package.

7. Three benefits of encapsulation:
	1) Clients cannot directly modify the object's variables.
	2) Hiding implementation details prevents clients from depending on things that might change. So we can impro
	ve our structure without any clients cautions.
	3) Encapsulation prevents clients from setting an object's variables arbitrary.

8. Methods that merely access/modify the internal values of a type are called getters/setters respectively.
	! WE NEEDN'T WRITE PREFIXES AS "Get"/"Fetch"/"Find"/"Lookup", BUT THE PREFIX "Set" IS NEEDED !

9. Encapsulation isn't always desireable.
*/
