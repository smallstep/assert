package assert

import (
	"errors"
	"reflect"
	"testing"
)

type tester struct {
	method string
	format string
	args   []interface{}
}

func tt() *tester {
	return &tester{}
}

func (t *tester) Errorf(format string, args ...interface{}) {
	t.method = "Errorf"
	t.format = format
	t.args = args
}

func (t *tester) Fatalf(format string, args ...interface{}) {
	t.method = "Fatalf"
	t.format = format
	t.args = args
}

func TestMessage(t *testing.T) {
	args := []interface{}{"%s", "a message"}
	if !reflect.DeepEqual(args, message(args, "default message")) {
		t.Fail()
	}
	if !reflect.DeepEqual([]interface{}{"default message"}, message(nil, "default message")) {
		t.Fail()
	}
}

func TestTrue(t *testing.T) {
	if True(tt(), false) {
		t.Fail()
	}
	if !True(tt(), true) {
		t.Fail()
	}
}

func TestFalse(t *testing.T) {
	if False(tt(), true) {
		t.Fail()
	}
	if !False(tt(), false) {
		t.Fail()
	}
}

func TestFatal(t *testing.T) {
	t1 := tt()
	Fatal(t1, false)
	if t1.method != "Fatalf" {
		t.Fail()
	}

	t2 := tt()
	Fatal(t2, true)
	if t2.method != "" {
		t.Fail()
	}
}

func TestFatalError(t *testing.T) {
	var err error
	t1 := tt()
	FatalError(t1, err)
	if t1.method != "" {
		t.Fail()
	}

	t2 := tt()
	FatalError(t2, errors.New("an error"))
	if t2.method != "Fatalf" {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	var err error
	if Error(tt(), err) {
		t.Fail()
	}
	if !Error(tt(), errors.New("an error")) {
		t.Fail()
	}
}

func TestNoError(t *testing.T) {
	if NoError(tt(), errors.New("an error")) {
		t.Fail()
	}
	var err error
	if !NoError(tt(), err) {
		t.Fail()
	}
}

func TestEquals(t *testing.T) {
	type myint int
	var nilPtr *int
	var myintPtr *myint
	var nilInterface interface{}
	var notNilInterface interface{}

	val, myval := 123, myint(123)
	notNilPtr := &val
	notNilInterface = val
	myintPtr = &myval

	type aType struct{}
	var nilType *aType
	notNilType := new(aType)

	tests := []struct {
		a, b interface{}
		res  bool
	}{
		{nil, nil, true},
		{0, 0, true},
		{0, nil, false},
		{nilPtr, nil, true},
		{notNilPtr, nil, false},
		{nilInterface, nil, true},
		{notNilInterface, nil, false},
		{nilPtr, nilPtr, true},
		{notNilPtr, notNilPtr, true},
		{nilInterface, nilInterface, true},
		{notNilInterface, notNilInterface, true},
		{nilPtr, nilInterface, true},
		{myint(123), myint(123), true},
		{nilType, nil, true},
		{nilType, nilType, true},
		{notNilType, nil, false},
		{notNilType, notNilType, true},
		{nil, nilPtr, true},
		{nil, nilInterface, true},
		{nil, nilType, true},
		{nilType, interface{}(nilType), true},
		{interface{}(notNilType), interface{}(notNilType), true},
		{interface{}(nilType), interface{}(nilType), true},
		{notNilType, interface{}(notNilType), true},
		{interface{}(notNilType), interface{}(notNilType), true},
		// not same type
		{123, myint(123), false},
		{notNilPtr, notNilInterface, false},
		{notNilPtr, myintPtr, false},
		{*notNilPtr, *myintPtr, false},
	}

	for i, tc := range tests {
		if Equals(tt(), tc.a, tc.b) != tc.res {
			t.Errorf("test %d with %v and %v failed", i, tc.a, tc.b)
		}
	}
}

func TestNil(t *testing.T) {
	var nilPtr *int
	var nilInterface interface{}
	var notNilInterface interface{}

	val := 123
	notNilPtr := &val
	notNilInterface = val

	type aType struct{}
	var nilType *aType
	notNilType := new(aType)

	tests := []struct {
		v   interface{}
		res bool
	}{
		{nil, true},
		{0, false},
		{1, false},
		{nilPtr, true},
		{notNilPtr, false},
		{nilInterface, true},
		{notNilInterface, false},
		{nilType, true},
		{notNilType, false},
	}

	for i, tc := range tests {
		if Nil(tt(), tc.v) != tc.res {
			t.Errorf("test %d with %v failed", i, tc.v)
		}
	}
}

func TestNotNil(t *testing.T) {
	var nilPtr *int
	var nilInterface interface{}
	var notNilInterface interface{}

	val := 123
	notNilPtr := &val
	notNilInterface = val

	type aType struct{}
	var nilType *aType
	notNilType := new(aType)

	tests := []struct {
		v   interface{}
		res bool
	}{
		{nil, false},
		{0, true},
		{1, true},
		{nilPtr, false},
		{notNilPtr, true},
		{nilInterface, false},
		{notNilInterface, true},
		{nilType, false},
		{notNilType, true},
	}

	for i, tc := range tests {
		if NotNil(tt(), tc.v) != tc.res {
			t.Errorf("test %d with %v failed", i, tc.v)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		v        interface{}
		expected int
		res      bool
	}{
		{nil, 0, false},
		{1234, 0, false},
		{"", 0, true},
		{"1234", 4, true},
		{[]int(nil), 0, true},
		{[]int{}, 0, true},
		{[]int{1, 2, 3}, 3, true},
		{[2]string{"foo", "bar"}, 2, true},
		{[...]string{"foo", "bar", "zar"}, 3, true},
		{map[string]int{}, 0, true},
		{map[string]int{"foo": 123}, 1, true},
		{make(chan int), 0, true},
	}

	for i, tc := range tests {
		if Len(tt(), tc.expected, tc.v) != tc.res {
			t.Errorf("test %d with %v failed", i, tc.v)
		}
	}
}

func TestPanic(t *testing.T) {
	withPanic := func() {
		panic("an error")
	}
	t1 := tt()
	if !Panic(t1, withPanic) || t1.method != "" {
		t.Fail()
	}
	t2 := tt()
	if Panic(t2, func() {}) || t2.method != "Errorf" {
		t.Fail()
	}
}

func TestType(t *testing.T) {
	type mytype string
	tests := []struct {
		e   interface{}
		v   interface{}
		res bool
	}{
		{0, 1, true},
		{0, "0", false},
		{mytype("a"), mytype("a"), true},
		{mytype("a"), mytype("b"), true},
		{mytype("a"), "a", false},
		{&tester{}, tt(), true},
		{tester{}, tt(), false},
	}

	for i, tc := range tests {
		if Type(tt(), tc.e, tc.v) != tc.res {
			t.Errorf("test %d with %v failed", i, tc.v)
		}
	}
}
