package assert

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

// Tester is an interface that testing.T implements, It has the methods used
// in the implementation of this package.
type Tester interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func reportError(t Tester, msg []interface{}) {
	_, file, line, _ := runtime.Caller(2)
	args := append([]interface{}{}, msg...)
	t.Errorf("\r\t%s:%d:\t%s", path.Base(file), line, fmt.Sprintln(args...))
}

func message(msg []interface{}, format string, a ...interface{}) []interface{} {
	if len(msg) > 0 {
		return msg
	}
	str := fmt.Sprintf(format, a...)
	return []interface{}{str}
}

// True checks that a condition is true.
func True(t Tester, condition bool, msg ...interface{}) bool {
	if !condition {
		msg = message(msg, "assert condition is not true")
		reportError(t, msg)
		return false
	}
	return true
}

// False checks that a condition is false.
func False(t Tester, condition bool, msg ...interface{}) bool {
	if condition {
		msg = message(msg, "assert condition is not false")
		reportError(t, msg)
		return false
	}
	return true
}

// Fatal checks that a condition is true or marks the test as failed and stop
// it's execution.
func Fatal(t Tester, condition bool, msg ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		args := append([]interface{}{}, msg...)
		t.Fatalf("\r\t%s:%d:\t%s", path.Base(file), line, fmt.Sprintln(args...))
	}
}

// FatalError checks that a error is nil or marks the test as failed and stop
// it's execution.
func FatalError(t Tester, err error, msg ...interface{}) {
	if err != nil {
		msg = message(msg, "error '%s' not expected", err)
		_, file, line, _ := runtime.Caller(1)
		args := append([]interface{}{}, msg...)
		t.Fatalf("\r\t%s:%d:\t%s", path.Base(file), line, fmt.Sprintln(args...))
	}
}

// Error checks if err is not nil.
func Error(t Tester, err error, msg ...interface{}) bool {
	if err == nil {
		msg = message(msg, "error expected but not found")
		reportError(t, msg)
		return false
	}
	return true
}

// NoError checks if err nil.
func NoError(t Tester, err error, msg ...interface{}) bool {
	if err != nil {
		msg = message(msg, "error '%s' not expected", err)
		reportError(t, msg)
		return false
	}
	return true
}

// Equals checks that expected and actual are equal.
func Equals(t Tester, expected, actual interface{}, msg ...interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}

	v1 := reflect.ValueOf(expected)
	v2 := reflect.ValueOf(actual)

	if isNilable(v1) && isNilable(v2) {
		switch {
		case v1.IsValid() && !v2.IsValid() && v1.IsNil():
			return true
		case !v1.IsValid() && v2.IsValid() && v2.IsNil():
			return true
		case v1.IsValid() && v2.IsValid() && v1.Type() == v2.Type() && v1.IsNil() && v2.IsNil():
			return true
		}
	}

	msg = message(msg, "'%v' and '%v' are not equal", expected, actual)
	reportError(t, msg)
	return false
}

// Nil checks that the value is nil.
func Nil(t Tester, value interface{}, msg ...interface{}) bool {
	var ret bool
	v := reflect.ValueOf(value)
	if isNilable(v) {
		ret = !v.IsValid() || v.IsNil()
	}

	if !ret {
		msg = message(msg, "nil expected and found '%v'", value)
		reportError(t, msg)
		return false
	}

	return true
}

// NotNil checks that the value is not nil.
func NotNil(t Tester, value interface{}, msg ...interface{}) bool {
	v := reflect.ValueOf(value)
	if isNilable(v) {
		if !v.IsValid() || v.IsNil() {
			msg = message(msg, "not nil expected and found '%v'", value)
			reportError(t, msg)
			return false
		}
	}

	return true
}

// Len checks that the application of len() to value match the expected value.
func Len(t Tester, expected int, value interface{}, msg ...interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() != expected {
			msg = message(msg, "len '%d' expected and found '%d'", expected, v.Len())
			reportError(t, msg)
			return false
		}
		return true
	default:
		msg = message(msg, "cannot apply built-in function len to '%s' (%v)", v.Kind(), value)
		reportError(t, msg)
		return false
	}
}

// Panic checks that the passed function panics.
func Panic(t Tester, f func(), msg ...interface{}) (ret bool) {
	defer func() {
		ret = true
		if r := recover(); r == nil {
			msg = message(msg, "function did not panic")
			reportError(t, msg)
			ret = false
		}
	}()
	f()
	return
}

// Type checks that the value matches the type of expected.
func Type(t Tester, expected, value interface{}, msg ...interface{}) bool {
	te := reflect.TypeOf(expected)
	tv := reflect.TypeOf(value)
	if te == tv {
		return true
	}
	msg = message(msg, "type '%T' expected and found '%T'", expected, value)
	reportError(t, msg)
	return false
}

// isNilable returns if the kind of v can be nil or not. It will return true
// for invalid values or if the kind is chan, func, interface, map, pointer,
// or slice; it will return false for the rest.
func isNilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// HasPrefix checks that the string contains the given prefix.
func HasPrefix(t Tester, s, p string, msg ...interface{}) bool {
	if strings.HasPrefix(s, p) {
		return true
	}
	msg = message(msg, "'%s' is not a prefix of '%s'", p, s)
	reportError(t, msg)
	return false
}
