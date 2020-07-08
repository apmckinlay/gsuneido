// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package hamcrest implements very basic hamcrest style asserts

For example:

	func TestStuff(t *testing.T) {
		Assert(t).That(2 * 4, Equals(6))
	}
*/
package hamcrest

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	gotesting "testing"
)

type testing interface {
	Error(err ...interface{})
}

// Asserter wraps a testing
type Asserter struct {
	t testing
}

// Assert returns an Asserter
func Assert(t testing) Asserter {
	return Asserter{t}
}

// True checks for true e.g. Assert(t).True(cond)
func (a Asserter) True(b bool) {
	if t, ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	if !b {
		a.Fail("expected true but it was false")
	}
}

// False checks for false e.g. Assert(t).False(cond)
func (a Asserter) False(b bool) {
	if t, ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	if b {
		a.Fail("expected false but it was true")
	}
}

// Tester is a function for That
type Tester func(interface{}) string

// That checks a value against a Tester
func (a Asserter) That(actual interface{}, test Tester) {
	if t, ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	err := test(actual)
	if err != "" {
		a.Fail(err)
	}
}

// Fail reports an error with its file and line
func (a Asserter) Fail(err string) {
	if t, ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	file, line := getLocation()
	a.t.Error(err + fmt.Sprintf(" {%s:%d}", file, line))
}

func getLocation() (string, int) {
	i := 1
	for ; i < 9; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok || strings.Contains(file, "testing/testing.go") {
			break
		}
	}
	_, file, line, ok := runtime.Caller(i - 1)
	if ok && i > 1 && i < 9 {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	return file, line
}

type Eq interface {
	Equal(interface{}) bool
}

// Equals returns a Tester
// that checks that the actual value is equal to the expected value
// float64's are compared as strings
// otherwise uses reflect.DeepEqual
func Equals(expected interface{}) Tester {
	return func(actual interface{}) string {
		if isNil(expected) && isNil(actual) {
			return ""
		}
		if a, ok := actual.(float64); ok {
			if e, ok := expected.(float64); ok {
				if strconv.FormatFloat(a, 'e', 15, 64) ==
					strconv.FormatFloat(e, 'e', 15, 64) {
					return ""
				}
			}
		}
		if intEqual(expected, actual) {
			return ""
		}
		if e, ok := expected.(Eq); ok {
			if e.Equal(actual) {
				return ""
			}
		} else if reflect.DeepEqual(expected, actual) {
			return ""
		}
		return fmt.Sprintf("\n    expect: %s\n    actual: %s\n    ",
			show(expected), show(actual))
	}
}

func isNil(x interface{}) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return v.IsNil()
	}
	return false
}

func intEqual(x interface{}, y interface{}) bool {
	var xi int64
	switch i := x.(type) {
	case int:
		xi = int64(i)
	case uint:
		xi = int64(i)
	case int8:
		xi = int64(i)
	case uint8:
		xi = int64(i)
	case int16:
		xi = int64(i)
	case uint16:
		xi = int64(i)
	case int32:
		xi = int64(i)
	case uint32:
		xi = int64(i)
	case int64:
		xi = int64(i)
	case uint64:
		xi = int64(i)
	default:
		return false
	}
	switch i := y.(type) {
	case int:
		return xi == int64(i)
	case uint:
		return xi == int64(i)
	case int8:
		return xi == int64(i)
	case uint8:
		return xi == int64(i)
	case int16:
		return xi == int64(i)
	case uint16:
		return xi == int64(i)
	case int32:
		return xi == int64(i)
	case uint32:
		return xi == int64(i)
	case int64:
		return xi == int64(i)
	case uint64:
		return xi == int64(i)
	default:
		return false
	}
}

func show(x interface{}) string {
	if _, ok := x.(string); ok {
		return fmt.Sprintf("%#v", x)
	}
	s1 := fmt.Sprintf("%v", x)
	s2 := fmt.Sprintf("%#v", x)
	if s1[0] == '[' {
		return s2
	}
	if s1 == s2 {
		return s1 + " (" + fmt.Sprintf("%T", x) + ")"
	}
	return s1 + " (" + s2 + ")"
}

// NotEquals returns a Tester for inequality using replect.DeepEqual
func NotEquals(expected interface{}) Tester {
	return func(actual interface{}) string {
		if !reflect.DeepEqual(expected, actual) {
			return ""
		}
		return fmt.Sprintf("expected %v not equal to %v but it was",
			expected, actual)
	}
}

// Like returns a Tester for comparing strings with whitespace standardized
func Like(expected interface{}) Tester {
	return func(actual interface{}) string {
		if like(expected.(string), actual.(string)) {
			return ""
		}
		return fmt.Sprintf("\nexpected: %s\nbut got: %s", expected, actual)
	}
}

func like(expected, actual string) bool {
	return canon(actual) == canon(expected)
}

func canon(s string) string {
	s = strings.TrimSpace(s)
	// can't use tr because it causes cycle in imports
	return rxlike.ReplaceAllString(s, " ")
}

var rxlike = regexp.MustCompile("[ \t\r\n]+")

// Panics returns a Tester that checks if a function panics
func Panics(expected string) Tester {
	return func(f interface{}) string {
		e := Catch(f.(func()))
		if e == nil {
			return fmt.Sprintf("expected panic of '%v' but it did not panic", expected)
		}
		if err, ok := e.(error); ok {
			e = err.Error()
		}
		if !strings.Contains(e.(string), expected) {
			return fmt.Sprintf("expected panic of '%v' but got '%v'", expected, e)
		}
		return ""
	}
}

// Catch calls the given function, catching and returning panics
func Catch(f func()) (result interface{}) {
	defer func() {
		if e := recover(); e != nil {
			//debug.PrintStack()
			result = e
		}
	}()
	f()
	return
}

// Comment decorates a Tester to add extra text to error messages
func (test Tester) Comment(items ...interface{}) Tester {
	return func(actual interface{}) string {
		err := test(actual)
		if err == "" {
			return ""
		}
		return err + " (" + fmt.Sprint(items...) + ")"
	}
}
