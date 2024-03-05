package params

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	httpTagKey    = "http"
	validTagKey   = "valid"
	emptyTagValue = ""
)

// search implements the /search url endpoint
func search(resp http.ResponseWriter, req *http.Request) {
	const (
		maxResultsNumber = 10
	)

	var data struct {
		Labels     []string `http:"1"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
		Weight     int      `http:"w"`
		Height     float64  `http:"h"`
	}

	data.MaxResults = maxResultsNumber

	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	// ...

	fmt.Fprintf(resp, "Search: %+v\n", data)
}

// Pack return a link based on the struct fields
func Pack(v any) (string, error) {

	const (
		maxBufferSize = 2048
		urlBody       = "localhost:8080/search?"
		equivalence   = '='
		urlAnd        = '&'
	)

	var (
		addrValue              = reflect.ValueOf(v)
		urlBuf    bytes.Buffer = *bytes.NewBuffer(make([]byte, 0, maxBufferSize))

		tagValue  string
		fieldKind reflect.Kind

		paramVal string
		err      error

		writeParam = func(paramVal string) {
			urlBuf.WriteString(tagValue)
			urlBuf.WriteByte(equivalence)
			urlBuf.WriteString(paramVal)
			urlBuf.WriteByte(urlAnd)
		}
	)

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return "", fmt.Errorf("unsupported type")
	}

	urlBuf.WriteString(urlBody)

	for i := 0; i < addrValue.NumField(); i++ {

		tagValue = getTagValue(addrValue.Type().Field(i), httpTagKey)
		fieldKind = addrValue.Field(i).Kind()

		switch fieldKind {
		case reflect.Slice:
			for j := 0; j < addrValue.Field(i).Len(); j++ {
				paramVal, err = paramValue(addrValue.Field(i).Index(j))
				if err != nil {
					return "", fmt.Errorf("getting field str representation: %s", err)
				}
				writeParam(paramVal)
			}

		default:
			paramVal, err = paramValue(addrValue.Field(i))
			if err != nil {
				return "", fmt.Errorf("getting field str representation: %s", err)
			}
			writeParam(paramVal)
		}
	}

	if urlBuf.Bytes()[urlBuf.Len()-1] == urlAnd {
		urlBuf.Truncate(urlBuf.Len() - 1)
	}

	return urlBuf.String(), nil
}

// getTagValue parses field's tag value and if it's presented
// returns this value in lowercase, otherwise returns the field'S
// name in lowercase
func getTagValue(field reflect.StructField, tagKey string) string {
	tag := field.Tag
	tagValue := tag.Get(tagKey)
	if tagValue == emptyTagValue {
		tagValue = strings.ToLower(field.Name)
	}

	return tagValue
}

// paramValue returns a string representation of the field value
// and the nil error value in none ones encountered,
// otherwise returns empty string and the corresponding non-nil error
func paramValue(reflectValue reflect.Value) (string, error) {

	var (
		strValue string
	)

	switch reflectValue.Kind() {
	case reflect.Invalid:
		return "", fmt.Errorf("invalid type")

	case reflect.String:
		if strValue := reflectValue.Interface().(string); strValue == "" {
			return "", fmt.Errorf("empty field value")
		} else {
			return strValue, nil
		}

	case reflect.Bool:
		value := reflectValue.Interface().(bool)
		if value {
			return "true", nil
		}
		return "false", nil

	case reflect.Int:
		value := reflectValue.Interface().(int)
		strValue = strconv.Itoa(value)
		return strValue, nil

	case reflect.Float64:
		value := reflectValue.Interface().(float64)
		strValue = strconv.FormatFloat(value, 'f', -1, 64)
		return strValue, nil

	default:
		return "", fmt.Errorf("unexpected reflection kind")
	}
}

// Unpack populates the fields of the struct pointed to by ptr
func Unpack(req *http.Request, ptr any) error {

	var (
		err     error
		addrVal reflect.Value

		fields map[string]reflect.Value

		httpTagValue string

		validTagValue    string
		queryLastElement string
	)

	if err = req.ParseForm(); err != nil {
		return err
	}
	fmt.Println(req.URL.Query())

	// Build a map of fields keyed by effective name.
	fields = make(map[string]reflect.Value)

	// Gets the struct addressable value
	addrVal = reflect.ValueOf(ptr).Elem()

	for i := 0; i < addrVal.NumField(); i++ {
		httpTagValue = getTagValue(addrVal.Type().Field(i), httpTagKey)
		validTagValue = getTagValue(addrVal.Type().Field(i), validTagKey)

		if _, ok := req.URL.Query()[validTagValue]; !ok {
			return fmt.Errorf("query doesn't have arguments for the parameter: %s", validTagValue)
		}

		queryLastElement = req.URL.Query()[validTagValue][len(req.URL.Query()[validTagValue])-1]
		if err = validateParamValue(validTagValue, queryLastElement); err != nil {
			return fmt.Errorf("validating query value: %s", err)
		}

		fields[httpTagValue] = addrVal.Field(i)

	}

	// Update struct field for each parameter in the request
	for name, values := range req.Form {
		f := fields[name]
		// Ignore unrecognized parameters
		if !f.IsValid() {
			continue
		}

		for _, value := range values {
			switch f.Kind() {
			case reflect.Slice:
				elem := reflect.New(f.Type().Elem()).Elem()
				if err = populate(elem, value); err != nil {
					return fmt.Errorf("%s %v", name, err)
				}
				f.Set(reflect.Append(f, elem))

			default:
				if err = populate(f, value); err != nil {
					return fmt.Errorf("%s %v", name, err)
				}
			}
		}
	}
	return nil
}

// populate takes field of a struct and value to be assigned to the field
// checks the validity of the value and if it's passed sets the value, returns the nil error,
// otherwise returns non-nil value
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

// validateParamValue takes the value to be assigned to the corresponding field, validation tag value if any,
// checks whether the passed value is correct and if everything is correct, returns nil error, otherwise returns
// non-nil value
func validateParamValue(validateTag, valueToValidate string) error {
	switch validateTag {
	case "email":
		const (
			invalidEmailLength = 0
			emailSeparator     = "@"

			leftSideInd         = 0
			validLeftSideLength = 5

			rightSideInd = 1
			validDomain  = "google.com"
		)

		email := valueToValidate

		if len(email) == invalidEmailLength {
			return fmt.Errorf("invalid email length: %d", len(email))
		}

		if !strings.Contains(email, emailSeparator) {
			return fmt.Errorf("password doesn't include %q: %d", '@', len(email))
		}

		emailPortions := strings.Split(email, emailSeparator)
		if len(emailPortions[leftSideInd]) < validLeftSideLength {
			return fmt.Errorf("email left side has invalid length: %d", len(emailPortions[leftSideInd]))
		}
		if emailPortions[rightSideInd] != validDomain {
			return fmt.Errorf("email domain isn't valid: %s", validDomain)
		}

		return nil
	case "zip":
		const (
			validZipLength = 6
			validZipPrefix = "100"
		)

		zip := valueToValidate

		if len(zip) < validZipLength {
			return fmt.Errorf("invalid zip length: %d", len(zip))
		}

		if !strings.HasPrefix(zip, validZipPrefix) {
			return fmt.Errorf("ivalid zip prefix: %s", zip)
		}

		return nil

	case "password":
		const (
			minPasswordLength = 8
			maxPasswordLength = 32
		)

		password := valueToValidate
		if !(len(password) >= minPasswordLength && len(password) <= maxPasswordLength) {
			return fmt.Errorf("invalid password length: %d", len(password))
		}

		return nil
	default:
		return fmt.Errorf("unsupported validity tag")
	}
}
