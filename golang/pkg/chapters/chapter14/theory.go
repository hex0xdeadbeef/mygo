package chapter14

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13. GENERICS

1. With generics we can declare and use functions or types that are written to work with any of a set of types provided by calling code.

2. We declared two functions with the different types because we must treat all of the differently (int64, float64).
 яч

3. With generics, we can write one fucntion instead of two. To support values of either type, that single function will need a way to declare what tyopes it supports. Calling code,
on the other hand, will need a way to specify whether it is calling an integer or float. To support this, we write a function that declares type parameters in addition to its ordinary
function parameters. These parameters make the function generic, enabling it to work with arguments of different types.

4. We can call a generic function with the ordinary function arguments.

5. Each type parameter has a type constraint that acts as a kind of meta-type for the type parameter. Each type constraint specifies the permissible type arguments that calling code can use
for the perspective type parameter. While a type parameter's constraint typically represents a set of types, at compile time the type parameter stands for a single type - the type provided
as a type argument by the calling code. if the type argument's type isn't allowed by the type parameter's constraint the code won't compile.
	1) Keep in mind that a type parameter must support all the operations the generic code is performing on it.

6. The declaration of a generic function is the following:
	func FuncName[type constraints] (a A, b B, vals []C, ...) A {
		...
	}

	1) The parameters A, B, C may be used within the function body as well. For example to declare a slice/map/array and so on.
	2) Type constraints are interface types

7. There's the specific predeclared interface "comparable". Interface that is implemented by all comparable types (booleans, numbers, strings, pointers, channels, arrays of comparable types,
structs whose fields are all comparable types). The comparable interface may only be used as a type parameter constraint, not as the type of a variable.

8. Using "|" specifies a union of the two types, meaning that this constraint allows either type. Either type will be permitted by the compiler as an argument in the calling code.

9. To call a generic function we should use the following statement:
	FuncName[type1, ..., typeN](params)

Also we can omit the type arguments in the fucntion call. Go can infer them from our code.

10. To run the code, in each call the compiler replaces the type parameters with the concrete types specified in that call, but this isn't always possible.
	1) If we need to call a generic function that had no arguments, we would need to include the type arguments in the function call.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.1 DECLARING A TYPE CONSTRAINT

1. We can declare a constraint as an interface. The constraint allows any type implementring the interface. For example:
	If we declare a type constraint interface with three methods, then use it with a type parameter in a generic function, type arguments used to call the function must have all of those
	methods.

2. Constraint interfaces can also refer to specific types using "|"
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.2 GENERIC TYPES

1. The declaration of a generic type is the following statement for example:
	type TypeName[K comparable, Y ..., ...] map[K]struct{} {
		...
	}
	When:
		1) In the square brackets we define a constraint type set we must use when creating some instances of it.

2. We can create the methods for generic types.d
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
13.2 CONSTRAINTS

1. The constraint types:
	1) [T any] - the "any" constraint accepts any type we want to pass, but it doesn't allow us to do much with the value because we don't really know what the value is. In this case
	we should provide the function with the logic that will define how the values will work with each other inside the function body. So the constraint "any" is barely useful.
	2) [C comparable]  - the "comparable" constraint says that any parameter of this type may be compared with "==" and "!=" operators. For example: the value may be any numeric type,
	string, bool, struct with all the comparable fields.
	3) [M MyInterface] - the "MyInterface" contraint may appear in the constraint list as well. We just create an interface and pass its name as the constraint. In this case when we put
	the interface name into the constraint set we should explicitly write out the preceding type parameter we declare in advance after the interface name. For example:
		func HasWithInterfaceConstraint[E Equalizer[E]](key E, values ...E) bool {
			for _, v := range values {
				if key.IsEqual(v) {
					return true
				}
			}

			return false
		}
	4) [UT ~underlyingType] - the "~underlyingType" constraint says that we can pass any arguments that have the same underlying type.
	5) [U type1|...|typeN] - the "type1|...|typeN" the type union means that inside we can use a type from this union. In this case we must check whether all the types support the operations we
	want to ptocess in the function body. If their do we won't get the compiler error, otherwise we will. For example:
		func Compare[T int | bool](a, b T) int {
			if a == b {
				return 0
			}

			if a < b { // invalid operation: a < b (type parameter T is not comparable with <)
				return -1
			}

			return 1
		}

	Also we can create the interfaces, that union the types inside the its body. The compiler will complain if any of the types don't support any operations. For example:
		type SimpleNumeric interface {
			int | int8 | int16 | int32 | int64 |
			uint | uint8 | uint16 | uint32 | uint64 |
			float32 | float64
		}

2. We can combine all the constraints in the interface type body. For example:

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
