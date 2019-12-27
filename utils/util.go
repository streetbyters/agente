package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3"
	html "html/template"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// HashPassword bcrypt hash generator with given password string and cost
func HashPassword(password string, cost int) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash)
}

// ComparePassword bcrypt compare with given hash password and raw password
func ComparePassword(hashPassword []byte, rawPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashPassword, rawPassword)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase string convert snake case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// InArray array search  with given value
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// ParseAndExecTemplateFromString parses string, creates template and exeutes, returns resulting string
func ParseAndExecTemplateFromString(s string, data interface{}) (string, error) {
	var buf bytes.Buffer

	t, err := template.New("test").Funcs(NewFuncMap()).Parse(s)
	if err != nil {
		return "", err
	}

	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// NewFuncMap template func generate
func NewFuncMap() map[string]interface{} {
	funcs := template.FuncMap{}
	funcs["iterate"] = func(start interface{}, count interface{}) []int64 {
		var i int64
		var Items []int64

		st := reflect.ValueOf(start)
		cn := reflect.ValueOf(count)

		for i = st.Int(); i <= st.Int()+cn.Int()-1; i++ {
			Items = append(Items, i)
		}
		return Items
	}

	funcs["add"] = add
	funcs["subtract"] = subtract
	funcs["multiply"] = multiply
	funcs["divide"] = divide
	funcs["iterate"] = iterate
	funcs["ToLower"] = strings.ToLower
	funcs["ToTitle"] = strings.ToTitle
	funcs["HasPrefix"] = strings.HasPrefix
	funcs["Contains"] = strings.Contains
	funcs["mustache"] = mustache
	funcs["attr"] = attr
	funcs["toText"] = toText
	funcs["stringInSlice"] = StringInSlice
	funcs["now"] = func() time.Time { return time.Now().Local() }

	funcs["Int64ToStr"] = Int64ToStr

	funcs["ieq"] = ieq
	funcs["startNewRow"] = startNewRow
	funcs["firstDateYear"] = func() time.Time {
		y, _, _ := time.Now().Date()
		return time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	funcs["lastDateYear"] = func() time.Time {
		y, _, _ := time.Now().Date()
		return time.Date(y, time.December, 31, 0, 0, 0, 0, time.UTC)
	}

	return funcs
}

// add returns the sum of a and b.
func add(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() + int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() + bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() + float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() + float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("add: unknown type for %q (%T)", av, a)
	}
}

// subtract returns the difference of b from a.
func subtract(a, b interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() - int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() - bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() - float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() - float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("subtract: unknown type for %q (%T)", av, a)
	}
}

// multiply returns the product of a and b.
func multiply(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() * int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() * bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() * float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() * float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() * bv.Float(), nil
		default:
			return nil, fmt.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("multiply: unknown type for %q (%T)", av, a)
	}
}

// divide returns the division of b from a.
func divide(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() / int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() / bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() / float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() / float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() / bv.Float(), nil
		default:
			return nil, fmt.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("divide: unknown type for %q (%T)", av, a)
	}
}

func iterate(start interface{}, count interface{}) []int64 {
	var i int64
	var Items []int64

	st := reflect.ValueOf(start)
	cn := reflect.ValueOf(count)

	for i = st.Int(); i <= st.Int()+cn.Int()-1; i++ {
		Items = append(Items, i)
	}
	return Items
}

func mustache(text string) html.HTML {
	return html.HTML("{{" + text + "}}")
}

func attr(text string) html.HTMLAttr {
	return html.HTMLAttr(text)
}

func toText(aval interface{}) interface{} {
	tmpVal := reflect.ValueOf(aval)

	switch tmpVal.Kind() {
	case reflect.Ptr:
		if aval != nil {
			return ""
		}

		return aval
	default:
		switch aval.(type) {
		case null.String:
			return aval.(null.String).String
		case null.Int:
			if aval.(null.Int).Valid {
				return Int64ToStr(aval.(null.Int).Int64)
			}
			return ""
		case null.Float:
			if aval.(null.Float).Valid {
				return FloatToStr(aval.(null.Float).Float64)
			}
			return ""
		case bool:
			if aval.(bool) {
				return `<i class="fa fa-check"></i>`
			}
			return ""
		default:
			return tmpVal
		}
	}

}

// Equality operator for interfaces which can have different types. For use in templates (ie options for select values)
func ieq(a interface{}, b interface{}) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	ak := av.Kind()
	bk := bv.Kind()

	if ak == reflect.Ptr {
		ak = reflect.Indirect(av).Kind()
		av = reflect.Indirect(av)
	}

	if bk == reflect.Ptr {
		bk = reflect.Indirect(bv).Kind()
		bv = reflect.Indirect(bv)
	}

	switch ak {
	case reflect.Int:
		val1 := av.Int()

		switch bk {
		case reflect.Int:
			return val1 == bv.Int()

		case reflect.Float64, reflect.Float32:
			return float64(val1) == bv.Float()

		case reflect.String:
			tmpval, err := strconv.Atoi(bv.String())
			if err != nil {
				return false
			}
			return val1 == int64(tmpval)
		}

	case reflect.Int64:
		val1 := av.Int()

		switch bk {
		case reflect.Int:
			return val1 == bv.Int()

		case reflect.Float64, reflect.Float32:
			return float64(val1) == bv.Float()

		case reflect.String:
			tmpval, err := strconv.Atoi(bv.String())
			if err != nil {
				return false
			}
			return val1 == int64(tmpval)
		}

	case reflect.String:
		val1 := av.String()

		switch bk {
		case reflect.Int, reflect.Int64:
			return val1 == Int64ToStr(bv.Int())

		case reflect.Float64, reflect.Float32:
			return val1 == strconv.FormatFloat(bv.Float(), 'f', -1, 64)

		case reflect.String:
			return val1 == bv.String()
		}

	}

	return a == b

}

func startNewRow(ndx int, colCount int) bool {
	if ndx == 0 {
		return true
	}

	rm := ndx % colCount

	if rm == 0 {
		return true
	}

	return false
}

// Int64ToStr is shorthand for strconv.FormatInt with base 10
func Int64ToStr(aval int64) string {
	res := strconv.FormatInt(aval, 10)
	return res
}

// StrToInt64 is shorthand for strconv.ParseInt with base 10, bitSize 64, returns 0 if parsing error occurs.
func StrToInt64(aval string) int64 {
	i, err := strconv.ParseInt(aval, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// StrToInt is shorhant for strconv.Atoi, returns 0 if parsing error occurs.
func StrToInt(aval string) int {
	i, err := strconv.Atoi(aval)
	if err != nil {
		return 0
	}
	return i
}

// FloatToStr formats float number for text representation. TODO: add formatting options as "#,##0.00"
func FloatToStr(aval float64) string {
	return fmt.Sprintf("%f", aval)
}

// StrToFloat is shorhand for strconv.ParseFşoat with bitSize 64, returns 0 if parsing error occurs.
func StrToFloat(aval string) float64 {
	i, err := strconv.ParseFloat(aval, 64)
	if err != nil {
		return 0
	}
	return i
}

// ParseTime parse string to time
func ParseTime(val string) (time.Time, error) {
	var err error

	if res, err := time.Parse("15:04", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("02.01.2006", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("01-02-2006", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("02.01.2006 15:04", val); err == nil {
		return res, nil
	}

	return time.Time{}, err
}

// StrToDate string to date
func StrToDate(aval string) (time.Time, error) {
	dt, err := time.Parse("02.01.2006", aval)
	if err != nil {
		return dt, err
	}

	return dt, nil
}

// StrToTime string to time
func StrToTime(aval string) (time.Time, error) {
	dt, err := time.Parse("15:04", aval)
	if err != nil {
		return dt, err
	}

	return dt, nil
}

// StrToTimeStamp string to timestamp
func StrToTimeStamp(aval string) (time.Time, error) {
	dt, err := time.Parse("02.01.2006 15:04", aval)
	if err != nil {
		return dt, err
	}

	return dt, nil
}

// JoinInt64Array int64 array join
func JoinInt64Array(lns []int64, sep string) string {
	lnsStr := make([]string, len(lns))
	for ndx, ln := range lns {
		lnsStr[ndx] = Int64ToStr(ln)
	}
	return strings.Join(lnsStr, sep)
}

// ParseInt string to itneger
func ParseInt(str string, base int, bitSize int) (i int64, flag bool) {
	i, err := strconv.ParseInt(str, base, bitSize)
	if err != nil {
		return i, true
	}
	return i, false
}

// StringInSlice string slice search
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Passkey generates a string passkey with an absolute length of 192.
func Passkey() string {
	var p []byte
	for i := 0; i < 9; i++ {
		b, _ := uuid.New().MarshalBinary()
		p = append(p, b...)
	}

	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(p)
}
