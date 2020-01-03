package assert

import (
	"fmt"
	"reflect"
	"strings"
)

// Tester is an interface that testing.T implements, It has the methods used
// in the implementation of this package.
type Tester interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
}

func reportError(t Tester, msg []interface{}) {
	args := append([]interface{}{}, msg...)
	t.Helper()
	t.Errorf(fmt.Sprintln(args...))
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
		t.Helper()
		reportError(t, message(msg, "assert condition is not true"))
		return false
	}
	return true
}

// False checks that a condition is false.
func False(t Tester, condition bool, msg ...interface{}) bool {
	if condition {
		t.Helper()
		reportError(t, message(msg, "assert condition is not false"))
		return false
	}
	return true
}

// Fatal checks that a condition is true or marks the test as failed and stop
// it's execution.
func Fatal(t Tester, condition bool, msg ...interface{}) {
	if !condition {
		msg = message(msg, "assert condition is not true")
		args := append([]interface{}{}, msg...)
		t.Helper()
		t.Fatalf(fmt.Sprintln(args...))
	}
}

// FatalError checks that a error is nil or marks the test as failed and stop
// it's execution.
func FatalError(t Tester, err error, msg ...interface{}) {
	if err != nil {
		msg = message(msg, "error '%s' not expected", err)
		args := append([]interface{}{}, msg...)
		t.Helper()
		t.Fatalf(fmt.Sprintln(args...))
	}
}

// Error checks if err is not nil.
func Error(t Tester, err error, msg ...interface{}) bool {
	if err == nil {
		t.Helper()
		reportError(t, message(msg, "error expected but not found"))
		return false
	}
	return true
}

// NoError checks if err nil.
func NoError(t Tester, err error, msg ...interface{}) bool {
	if err != nil {
		t.Helper()
		reportError(t, message(msg, "error '%s' not expected", err))
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

	t.Helper()
	reportError(t, message(msg, "'%v' and '%v' are not equal", expected, actual))
	return false
}

// NotEquals checks that expected and actual are not equal.
func NotEquals(t Tester, expected, actual interface{}, msg ...interface{}) bool {
	v1 := reflect.ValueOf(expected)
	v2 := reflect.ValueOf(actual)

	if isNilable(v1) && isNilable(v2) {
		switch {
		case v1.IsValid() && !v2.IsValid() && v1.IsNil():
			fallthrough
		case !v1.IsValid() && v2.IsValid() && v2.IsNil():
			fallthrough
		case v1.IsValid() && v2.IsValid() && v1.Type() == v2.Type() && v1.IsNil() && v2.IsNil():
			t.Helper()
			reportError(t, message(msg, "'%v' and '%v' are equal", expected, actual))
			return false
		}
	}

	if reflect.DeepEqual(expected, actual) {
		t.Helper()
		reportError(t, message(msg, "'%v' and '%v' are equal", expected, actual))
		return false
	}

	return true
}

// Nil checks that the value is nil.
func Nil(t Tester, value interface{}, msg ...interface{}) bool {
	var ret bool
	v := reflect.ValueOf(value)
	if isNilable(v) {
		ret = !v.IsValid() || v.IsNil()
	}

	if !ret {
		t.Helper()
		reportError(t, message(msg, "nil expected and found '%v'", value))
		return false
	}

	return true
}

// NotNil checks that the value is not nil.
func NotNil(t Tester, value interface{}, msg ...interface{}) bool {
	v := reflect.ValueOf(value)
	if isNilable(v) {
		if !v.IsValid() || v.IsNil() {
			t.Helper()
			reportError(t, message(msg, "not nil expected and found '%v'", value))
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
			t.Helper()
			reportError(t, message(msg, "len '%d' expected and found '%d'", expected, v.Len()))
			return false
		}
		return true
	default:
		t.Helper()
		reportError(t, message(msg, "cannot apply built-in function len to '%s' (%v)", v.Kind(), value))
		return false
	}
}

// Panic checks that the passed function panics.
func Panic(t Tester, f func(), msg ...interface{}) (ret bool) {
	t.Helper()
	defer func() {
		ret = true
		if r := recover(); r == nil {
			t.Helper()
			reportError(t, message(msg, "function did not panic"))
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
	t.Helper()
	reportError(t, message(msg, "type '%T' expected and found '%T'", expected, value))
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
	t.Helper()
	reportError(t, message(msg, "'%s' is not a prefix of '%s'", p, s))
	return false
}
