package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
