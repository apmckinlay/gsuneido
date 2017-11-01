/*
Package hamcrest implements very basic hamcrest style asserts

For example:

	func TestStuff(t *testing.T) {
		Assert(t).That(2 * 4, Equals(6))
	}
*/
package hamcrest

import "fmt"
import (
	"reflect"
	"regexp"
)
import "runtime"
import "strings"
import gotesting "testing"

type testing interface {
	Error(err ...interface{})
}

type Asserter struct {
	t testing
}

func Assert(t testing) Asserter {
	return Asserter{t}
}

type Tester func(interface{}) string

func (a Asserter) True(b bool) {
	if b != true {
		a.Fail("expected true but it was false")
	}
}

func (a Asserter) False(b bool) {
	if b != false {
		a.Fail("expected false but it was true")
	}
}

func (a Asserter) That(actual interface{}, test Tester) {
	if t,ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	err := test(actual)
	if err != "" {
		a.Fail(err)
	}
}

func (a Asserter) Fail(err string) {
	if t,ok := a.t.(*gotesting.T); ok {
		t.Helper() // skip this function when printing file/line info
	}
	file, line := getLocation()
	a.t.Error(err + fmt.Sprintf(" {%s:%d}", file, line))
}

func getLocation() (file string, line int) {
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

// Equals checks that the actual value is equal to the expected value
// using reflect.DeepEquals
func Equals(expected interface{}) Tester {
	return func(actual interface{}) string {
		if reflect.DeepEqual(expected, actual) {
			return ""
		}
		return fmt.Sprintf("expected: %#v but got: %#v", expected, actual)
	}
}

func NotEquals(expected interface{}) Tester {
	return func(actual interface{}) string {
		if !reflect.DeepEqual(expected, actual) {
			return ""
		}
		return fmt.Sprintf("expected %v not equal to %v but it was",
			expected, actual)
	}
}

func Like(expected interface{}) Tester {
	return func(actual interface{}) string {
		if like(expected.(string), actual.(string)) {
			return ""
		}
		return fmt.Sprintf("expected: %s\nbut got: %s", expected, actual)
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

type runnable func()

func Panics(expected string) Tester {
	return func(f interface{}) (result string) {
		defer func() {
			if e := recover(); e != nil {
				if strings.Contains(e.(string), expected) {
					result = ""
				} else {
					result = fmt.Sprintf("expected panic of '%v' but got '%v'",
						expected, e)
				}
			}
		}()
		f.(func())()
		return fmt.Sprintf("expected panic of '%v' but it did not panic", expected)
	}
}

// Comment decorates a Tester to add extra text to error messages
func (test Tester) Comment(items ...interface{}) Tester {
	return func(actual interface{}) string {
		err := test(actual)
		if err == "" {
			return ""
		} else {
			return err + " (" + fmt.Sprint(items...) + ")"
		}
	}
}
