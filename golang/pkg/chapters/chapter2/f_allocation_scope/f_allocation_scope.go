package f_allocation_scope

// import is executed by compiler before all package
import (
	"fmt"
	"log"
	r "math/rand" // This is the alias
	"os"
)

// Allocation on heap/stack doesn't depend on using var/new(...)
// It depends whether variable is reachable or not
var global *int

func f() {
	var x int // It's heap-allocated 'cause it's still reachable from the variable global
	x = 1
	global = &x
}

func g() {
	var y = new(int) // It's stack allocated 'case it won't be reacable after g() has returned
	*y = 10
}

func ScopeShadowingFirst() {
	x := "hello" // "function" scope x, it'll be shadowed by for scope x variable
	for i := 0; i < len(x); i++ {
		x := x[i] // "for" scope x, it'll be shadowed by if scope x
		if x != '!' {
			x := x + 'A' - 'a'  // "if" scope x
			fmt.Printf("%c", x) // this statement uses the last mentioned x (x := x + 'A' - 'a' // "if" scope x)
		}
	}
}

func ScopeShadowingSecond() {
	x := "hello"          // "function" scope x, it'll be shadowed by the "for" scape implicit variable
	for _, x := range x { // The value that is fetched from string x is the new variable
		x := x + 'A' - 'a'  // This x is the new variable that is assigned with expression
		fmt.Printf("%c", x) // This statement uses the last mentioned x (x := x + 'A' - 'a' // "if" scope x)
	}
}

func xRet() int {
	if value := r.Intn(100); value > 50 {
		return 1
	} else {
		return 0
	}
}

func yRet(x int) int {
	if value := r.Intn(100); value > 50 {
		return 0
	} else {
		return 1
	}
}

func ScopeShadowingThird() {
	if x := xRet(); x == 0 { // 1. Assignment x a value from xRet(). 2. Evaluating the expression
		fmt.Println(x) // Printing
	} else if y := yRet(x); y == x { /* 1. Assignment y a value from yRet(). 2. Evaluating the expression by using
		the evaluated value of x*/
		fmt.Println(x, y)
	} else {
		fmt.Println(x, y) // Printing both variables by by using the evaluated values of x and y
	}
	// fmt.Println(x, y) // Here's no way to reach both x and y variables 'cause they are out of general if statement
}

func ScopeShadowingFourth(fileNames ...string) {
	if file, err := os.Create(fileNames[0]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		defer file.Close()
		var fileData []byte
		_, err := file.Read(fileData)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		os.Stdout.Write(fileData)
	}
	// defer file.Close() It will not be accessible 'cause file variable scope has ended above
	// f.ReadByte() -||-

}

var cwd string

func DefeatGlobalVarShadowing() {

	fmt.Println("BEFORE | The address in memory of cwd:", &cwd)
	/* Here's the local cwd shades the upper declared global cwd */
	// cwd, err := os.Getwd()

	/* Here's a method that defeats this problem */
	var err error
	cwd, err = os.Getwd() // cwd has assigned, but hasn't declared. It's

	// cwd, err := os.Getwd() // Here's the local cwd shades the upper declared global cwd
	if err != nil {
		log.Fatalf("os.Getwd failed:", err)
	}

	fmt.Println("AFTER | The address in memory of cwd:", &cwd, "| The working directory is:", cwd)
}
