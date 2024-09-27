package chapter7

import (
	"bufio"
	"errors"
	"fmt"
	. "golang/pkg/chapters/chapter3/b_floats/threedsurface"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// The filler for a variables presence map
type Empty struct{}

var depth = -1

// The variables of this type will contain any expression
type Expression interface {
	// Returns the value of this Expr in the environment env.
	Eval(env Environment) float64
	// Reports errors in this Expr and adds its Vars to the set
	Check(vars map[Variable]Empty) error
	// Adds the string representation on a corresponding tree level
	String(levels map[int]string)
}

// Identifies a variable e.g. x
type Variable string

/*
Just returns a respective float64 value from the environment if the corresponding value of the key is presented in environment
*/
func (v Variable) Eval(env Environment) float64 {
	return env[v]
}

/*
1) Put the variable name into the map of the variables
2) Return nil
*/
func (v Variable) Check(vars map[Variable]Empty) error {
	vars[v] = Empty{}
	return nil
}

func (v Variable) String(levels map[int]string) {
	depth++
	_, doesLevelPresented := levels[depth]
	if !doesLevelPresented {
		levels[depth] = fmt.Sprintf("%d| ", depth)
	}
	levels[depth] += (string(v) + ", ")

	depth--
}

// Numeric constant e.g. 3.141
type literal float64

/*
Just returns the respective float64 value
*/
func (l literal) Eval(_ Environment) float64 {
	return float64(l)
}

/*
Just return nil
*/
func (l literal) Check(vars map[Variable]Empty) error {
	return nil
}

func (l literal) String(levels map[int]string) {
	depth++
	_, doesLevelPresented := levels[depth]
	if !doesLevelPresented {
		levels[depth] = fmt.Sprintf("%d| ", depth)
	}
	levels[depth] += fmt.Sprintf("%g, ", l)
	depth--
}

// Represents a unaryOp operator expression e.g. -x
type unaryOp struct {
	operationCharacter rune       // "-" or "+"
	operand            Expression // any expression
}

/*
1) Choose an appropriate switch statement
2) Evaluate the right side expression
3) Apply the left side operation "+" or "-"
4) Return the result
*/
func (u unaryOp) Eval(env Environment) float64 {
	switch u.operationCharacter {
	case '+':
		return +u.operand.Eval(env)
	case '-':
		return -u.operand.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported unary operator: %q", u.operationCharacter))
	}
}

/*
1) Check the validity of an operation given
2) If everything is a-okay, return nil, otherwise return the error
*/
func (u unaryOp) Check(vars map[Variable]Empty) error {
	// Check whether the operation is valid
	if !strings.ContainsRune("+-", u.operationCharacter) {
		return fmt.Errorf("unexpected unary op %q", u.operationCharacter)
	}
	return u.operand.Check(vars)

}

func (u unaryOp) String(levels map[int]string) {
	depth++
	_, doesLevelPresented := levels[depth]
	if !doesLevelPresented {
		levels[depth] = fmt.Sprintf("%d| ", depth)
	}
	levels[depth] += fmt.Sprintf("%cU, ", u.operationCharacter)
	u.operand.String(levels)
	depth--
}

// Represents a binaryOp operator expression e.g. x+y
type binaryOp struct {
	operationCharacter        rune // "+" or "-" or "*" or "/"
	leftOperand, rightOperand Expression
}

/*
1) Chose an appropriate switch statement
2) Evaluate the left operand value
3) Evaluate the right operand value
4) Return the result of operation applied to the operands
*/
func (b binaryOp) Eval(env Environment) float64 {
	switch b.operationCharacter {
	case '+':
		return b.leftOperand.Eval(env) + b.rightOperand.Eval(env)
	case '-':
		return b.leftOperand.Eval(env) - b.rightOperand.Eval(env)
	case '*':
		return b.leftOperand.Eval(env) * b.rightOperand.Eval(env)
	case '/':
		return b.leftOperand.Eval(env) / b.rightOperand.Eval(env)
	default:
		panic(fmt.Sprintf("unsupported binary operator: %q", b.operationCharacter))
	}
}

/*
1) Check the validity of an operation given
2) If everything is a-okay, return nil, otherwise return the error
*/
func (b binaryOp) Check(vars map[Variable]Empty) error {
	// Check whether an operation is valid
	if !strings.ContainsRune("+-*/", b.operationCharacter) {
		return fmt.Errorf("unexpected binary op %q", b.operationCharacter)
	}
	if err := b.leftOperand.Check(vars); err != nil {
		return err
	}
	return b.rightOperand.Check(vars)
}

func (b binaryOp) String(levels map[int]string) {
	depth++
	_, doesLevelPresented := levels[depth]
	if !doesLevelPresented {
		levels[depth] = fmt.Sprintf("%d| ", depth)
	}
	levels[depth] += fmt.Sprintf("B %c B, ", b.operationCharacter)
	b.leftOperand.String(levels)
	b.rightOperand.String(levels)
	depth--
}

// Represents a function functionCall expression e.g. sin(x)
type functionCall struct {
	functionName string // "pow" or "sqrt" or "sin"
	arguments    []Expression
}

/*
1) Choose an appropriate switch statement
2) Evaluate the argument expression
3) Call the function with an evaluated argument
4) Return the result of a function call
*/
func (fc functionCall) Eval(env Environment) float64 {
	switch fc.functionName {
	case "pow":
		return math.Pow(fc.arguments[0].Eval(env), fc.arguments[1].Eval(env))
	case "sin":
		return math.Sin(fc.arguments[0].Eval(env))
	case "sqrt":
		return math.Sqrt(fc.arguments[0].Eval(env))
	default:
		panic(fmt.Sprintf("unsupported function call: %q", fc.functionName))
	}
}

// The map to check an arity of a function
var paramsNum = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

/*
1) Check the presence of the function gievn in the function map
2) Check whether params count is valid or not
3) Apply the check for each argument of function
4) If all checks are successful, return nil, otherwise return the corresponding value.
*/
func (fc functionCall) Check(vars map[Variable]Empty) error {

	arity, ok := paramsNum[fc.functionName]

	if !ok {
		return fmt.Errorf("unknown function %q", fc.functionName)
	}

	if len(fc.arguments) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", fc.functionName, len(fc.arguments), arity)
	}

	for _, arg := range fc.arguments {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}

	return nil
}

func (fc functionCall) String(levels map[int]string) {
	depth++
	_, doesLevelPresented := levels[depth]
	if !doesLevelPresented {
		levels[depth] = fmt.Sprintf("%d| ", depth)
	}
	str := fc.functionName + "( "
	for i := 0; i < len(fc.arguments); i++ {
		str += fmt.Sprintf("arg%d ", i+1)
	}
	str += ")"
	levels[depth] += (str + ", ")
	for _, arg := range fc.arguments {
		arg.String(levels)
	}
	depth--
}

// The mapping variables to values
type Environment map[Variable]float64

func TestEval() {
	tests := []struct {
		expr string
		env  Environment
		want string
	}{
		{"sqrt(A/pi)", Environment{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x,3) + pow(y,3)", Environment{"x": 12, "y": 1}, "1729"},
		{"5/9*(F-32)", Environment{"F": -40}, "-40"},
		{"5/9*(F-32)", Environment{"F": 32}, "0"},
		{"5/9*(F-32)", Environment{"F": 212}, "100"},
	}

	var prevExpr string

	for _, test := range tests {
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}

		expr, err := Parse(test.expr)
		if err != nil {
			// t.Error(err)
			continue
		}

		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			fmt.Println(fmt.Errorf("%s.Eval() in %v = %q, want %q\n", test.expr, test.env, got, test.want))
		}
	}
}

/*
Performs all the checks of input, expression, variables.
*/
func parseAndCheck(input string) (Expression, error) {
	// Checks whether the string is empty
	if input == "" {
		return nil, fmt.Errorf("empty expression")
	}

	// Checks whether there are any lexical errors
	expr, err := Parse(input)
	if err != nil {
		return nil, err
	}

	// Checks whether there are any semantic errors
	vars := make(map[Variable]Empty)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}

	// Check the validity of the every variable
	for variable := range vars {
		if variable != "x" && variable != "y" && variable != "r" {
			return nil, fmt.Errorf("undefined variable %s", variable)
		}
	}

	return expr, nil
}

func plot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get an expression from the user input and perform its validation check
	expr, err := parseAndCheck(r.Form.Get("expr"))
	if err != nil {
		http.Error(w, fmt.Sprintf("bad expr: %s", err), http.StatusBadRequest)
		return
	}

	// Plot a surface and send to user the image of the plotted surface
	GetSurface(w,
		func(x, y float64) float64 {
			r := math.Hypot(x, y) // distance from (0,0)
			return expr.Eval(Environment{"x": x, "y": y, "r": r})
		})
}

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/plot", plot)

	http.ListenAndServe("localhost:8080", mux)
}

/* HOMEWORK */

// 7.13
func GetStringRepresentation(strExpression string) error {

	expr, err := Parse(strExpression)
	if err != nil {
		return err
	}

	// Declare the necessary variables to get string representation
	levels := make(map[int]string, 0)

	// Put the tree into the map
	expr.String(levels)

	// Print all the levels of the tree
	for i := 0; i < len(levels); i++ {
		fmt.Println(levels[i])
	}

	return nil
}

// 7.15
func EvaluateInputExpression() error {
	// Get string expression from user
	strExpr, err := getExpressionFromUser()
	if err != nil {
		return fmt.Errorf("while getting input from user; %s", err)
	}

	// Validate string expression, get an expression and its variables map
	expr, exprVars, err := validateExpression(strExpr)
	if err != nil {
		return fmt.Errorf("validating expression string \"%s\"; %s", strExpr, err)
	}

	var (
		argumentsValues []string
		environment     map[Variable]float64
	)
	// If the variables map length isn't zero, get the variables values
	if len(exprVars) > 0 {
		if argumentsValues, err = getArgumentsValues(exprVars); err != nil {
			return fmt.Errorf("getting arguments values from a user for params %s; %s", getParamsString(exprVars), err)
		}
		if environment, err = getArgumentsEnvironment(exprVars, argumentsValues); err != nil {
			return fmt.Errorf("creating the environment for %s; %s", getParamsString(exprVars), err)
		}
	}

	// Evaluate the expression
	result := expr.Eval(environment)

	fmt.Printf("-------------\nExpression -> %s\nParams -> %v\nResult -> %g\n-------------\n", strExpr, environment, result)

	return nil
}

/*
Offers a user to write a math expression, tracks the input attempts count. If attempts haven't been exceeded returns the string
expression and non-nil error, otherwise an empty string and the corresponding error.
*/
func getExpressionFromUser() (input string, err error) {
	const maxAttemptsCount = 5

	var (
		userInput    string
		attemptCount int
	)

	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Fprintf(os.Stdout, "Write an expression -> ")

	// Get an expression from a user
	for scanner.Scan() {
		userInput = scanner.Text()
		attemptCount++

		if attemptCount == maxAttemptsCount {
			return "", errors.New("Attempts limit exceeded")
		}

		if len(strings.Trim(userInput, "\t\n\v\f\r \u0085\u00A0")) == 0 {
			fmt.Printf("invalid expression; expr: \"%s\"\n", userInput)
			fmt.Printf("Write an expression =>")
			continue
		}

		fmt.Printf("The expression \"%s\" has been initially accepted.\n", userInput)
		break
	}

	return userInput, nil
}

/*
Validates the string expression. If it encountered no errors, returns the given expression, its variables map, nil error, otherwise
the corresponding encountered error.
*/
func validateExpression(input string) (Expression, map[Variable]Empty, error) {
	// Checks whether the string is empty
	if input == "" {
		return nil, nil, fmt.Errorf("empty expression")
	}

	// Checks whether there are any lexical errors
	expr, err := Parse(input)
	if err != nil {
		return nil, nil, err
	}

	// Checks whether there are any semantic errors
	vars := make(map[Variable]Empty)
	if err := expr.Check(vars); err != nil {
		return nil, nil, err
	}

	return expr, vars, nil
}

/*
Offers a user to write an values, tracks the input attempts count. If attempts haven't been exceeded returns the values map and non-nil
error, otherwise a nil map and the corresponding error.
*/
func getArgumentsValues(variables map[Variable]Empty) ([]string, error) {
	const maxAttemptsCount = 5

	var (
		userInput    []string
		attemptCount int
	)

	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Printf("Write values for %s params separated by commas -> ", getParamsString(variables))

	// Get values for parameters from user
	for scanner.Scan() {
		userInput = strings.Split(strings.Trim(scanner.Text(), "\t\n\v\f\r \u0085\u00A0"), ",")

		attemptCount++

		if attemptCount == maxAttemptsCount {
			return nil, errors.New("Attempts limit exceeded")
		}

		if len(userInput) != len(variables) {
			fmt.Printf("invalid values; values: \"%v\"\n", userInput)
			fmt.Printf("Write values for %s params separated by commas -> ", getParamsString(variables))
			continue
		}

		fmt.Printf("The values \"%v\" has been initially accepted.\n", userInput)
		break
	}

	return userInput, nil
}

/*
Converts required params list into a string
*/
func getParamsString(variables map[Variable]Empty) string {
	result := "| "
	for key := range variables {
		result += (string(key) + " ")
	}
	result += "|"

	return result
}

/*
Makes a map of argument/value pairs.
*/
func getArgumentsEnvironment(variablesMap map[Variable]Empty, values []string) (map[Variable]float64, error) {
	environment := make(map[Variable]float64)
	var ind int

	for key := range variablesMap {
		parsedValue, err := strconv.ParseFloat(values[ind], 64)
		if err != nil {
			return nil, fmt.Errorf("%s; value: \"%s\"", err, values[ind])
		}
		environment[key] = parsedValue

		ind++
	}

	return environment, nil
}
