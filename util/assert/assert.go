// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package assert helps writing assertions for tests.
// Benefits include brevity, clarity, and helpful error messages.
//
// If .T(t) is specified, failures are reported with t.Error
// which means multiple errors may be reported.
// If .T(t) is not specified, panic is called with the error string.
//
// For example:
//
//	assert.That(a || b)
//	assert.This(x).Is(y)
//	assert.T(t).This(x).Like(y)
//	assert.Msg("first time").That(a || b)
//	assert.T(t).Msg("second").This(fn).Panics("illegal")
//
// Use a variable to avoid specifying .T(t) repeatedly:
//
//	assert := assert.T(t)
//	assert.This(x).Is(y)
package assert

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dbg"
)

type assert struct {
	t   *testing.T
	msg []any
}

// T specifies a *testing.T to use for reporting errors
func T(t *testing.T) assert {
	return assert{t: t}
}

// Msg adds additional information to print with the error.
// It can be useful with That/True/False where the error is not very helpful.
func Msg(args ...any) assert {
	return assert{msg: args}
}

// Msg adds additional information to print with the error.
// It can be useful with That/True/False where the error is not very helpful.
func (a assert) Msg(args ...any) assert {
	a.msg = args
	return a
}

// Msg adds additional information to print with the error.
// It can be useful with That/True/False where the error is not very helpful.
func (v value) Msg(args ...any) value {
	v.assert.msg = args
	return v
}

// Nil gives an error if the value is not nil.
// It handles nil func/pointer/slice/map/channel using reflect.IsNil
// For performance critical code, consider using That(value == nil)
func Nil(v any) {
	assert{}.This(v).Is(nil)
}

// Nil gives an error if the value is not nil.
// It handles nil func/pointer/slice/map/channel using reflect.IsNil
// For performance critical code, consider using That(value == nil)
func (a assert) Nil(v any) {
	if a.t != nil {
		a.t.Helper()
	}
	a.This(v).Is(nil)
}

// True gives an error if the value is not true.
// True(x) is the same as That(x)
func True(b bool) {
	assert{}.That(b)
}

// True gives an error if the value is not true.
// True(x) is the same as That(x)
func (a assert) True(b bool) {
	if a.t != nil {
		a.t.Helper()
	}
	a.That(b)
}

// False gives an error if the value is not true.
// False(x) is the same as That(!x)
func False(b bool) {
	assert{}.That(!b)
}

// False gives an error if the value is not true.
// False(x) is the same as That(!x)
func (a assert) False(b bool) {
	if a.t != nil {
		a.t.Helper()
	}
	a.That(!b)
}

// That gives an error if the value is not true.
// That(x) is the same as True(x)
func That(cond bool) {
	if !cond {
		panic("ASSERT FAILED")
	}
}

// That gives an error if the value is not true.
// That(x) is the same as True(x)
func (a assert) That(cond bool) {
	if !cond {
		if a.t != nil {
			a.t.Helper()
		}
		a.fail()
	}
}

// This sets a value to be compared e.g. with Is or Like
func This(v any) value {
	return value{value: v}
}

// This sets a value to be compared e.g. with Is or Like.
// It is usually the actual value and Is gives the expected.
func (a assert) This(v any) value {
	return value{assert: a, value: v}
}

type value struct {
	value  any
	assert assert
}

// Is gives an error if the given expected value is not the same
// as the actual value supplied to This.
// Accepts as equivalent: different nils and different int types.
// Compares floats via string forms.
// Uses an Equal method if available on the expected value.
// Finally, uses reflect.DeepEqual.
func (v value) Is(expected any) {
	if !Is(v.value, expected) {
		if v.assert.t != nil {
			v.assert.t.Helper()
		}
		v.assert.fail("expected: ", show(expected),
			"\nactual: ", show(v.value))
	}
}

// Isnt gives an error if the given expected value is the same
// as the actual value supplied to This.
func (v value) Isnt(expected any) {
	if Is(v.value, expected) {
		if v.assert.t != nil {
			v.assert.t.Helper()
		}
		v.assert.fail("expected not: ", show(expected), " but it was")
	}
}

func Is(actual, expected any) bool {
	if isNil(expected) && isNil(actual) {
		return true
	}
	if a, ok := actual.(float64); ok {
		if e, ok := expected.(float64); ok {
			if strconv.FormatFloat(a, 'e', 15, 64) ==
				strconv.FormatFloat(e, 'e', 15, 64) {
				return true
			}
		}
	}
	if intEqual(expected, actual) {
		return true
	}
	type equal interface {
		Equal(any) bool
	}
	if e, ok := expected.(equal); ok {
		if e.Equal(actual) {
			return true
		}
	} else if reflect.DeepEqual(expected, actual) {
		return true
	}
	return false
}

func isNil(x any) bool {
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

func intEqual(x any, y any) bool {
	var x64 int64
	switch x := x.(type) {
	case int:
		x64 = int64(x)
	case uint:
		x64 = int64(x)
	case int8:
		x64 = int64(x)
	case uint8:
		x64 = int64(x)
	case int16:
		x64 = int64(x)
	case uint16:
		x64 = int64(x)
	case int32:
		x64 = int64(x)
	case uint32:
		x64 = int64(x)
	case int64:
		x64 = int64(x)
	case uint64:
		x64 = int64(x)
	default:
		return false
	}
	switch y := y.(type) {
	case int:
		return x64 == int64(y)
	case uint:
		return x64 == int64(y)
	case int8:
		return x64 == int64(y)
	case uint8:
		return x64 == int64(y)
	case int16:
		return x64 == int64(y)
	case uint16:
		return x64 == int64(y)
	case int32:
		return x64 == int64(y)
	case uint32:
		return x64 == int64(y)
	case int64:
		return x64 == int64(y)
	case uint64:
		return x64 == int64(y)
	default:
		return false
	}
}

func show(x any) string {
	if _, ok := x.(string); ok {
		return fmt.Sprintf("%#v", x)
	}
	if r, ok := x.(rune); ok {
		return "'" + string(r) + "'"
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

// Like compares strings with whitespace standardized.
// Leading and trailing whitespace is removed,
// runs of whitespace are converted to a single space.
func (v value) Like(expected any) {
	exp := expected.(string)
	val := v.value.(string)
	if !like(exp, val) {
		if v.assert.t != nil {
			v.assert.t.Helper()
		}
		sep := " "
		if strings.Contains(exp, "\n") || strings.Contains(val, "\n") {
			sep = "\n"
		}
		v.assert.fail("expected:" + sep + exp + "\nbut got:" + sep + val)
	}
}

func like(expected, actual string) bool {
	return canon(actual) == canon(expected)
}

func canon(s string) string {
	s = strings.TrimSpace(s)
	s = leadingSpace.ReplaceAllString(s, "")
	s = trailingSpace.ReplaceAllString(s, "")
	s = whitespace.ReplaceAllString(s, " ")
	return s
}

var leadingSpace = regexp.MustCompile("(?m)^[ \t]+")
var trailingSpace = regexp.MustCompile("(?m)[ \t]+$")
var whitespace = regexp.MustCompile("[ \t]+")

// Panics checks if a function panics
func (v value) Panics(expected string) {
	e := Catch(v.value.(func()))
	if e == nil {
		if v.assert.t != nil {
			v.assert.t.Helper()
		}
		v.assert.fail(fmt.Sprintf("expected panic with '%v' but it did not panic",
			expected))
		return
	}
	if err, ok := e.(error); ok {
		e = err.Error()
	}
	if !strings.Contains(e.(string), expected) {
		if v.assert.t != nil {
			v.assert.t.Helper()
		}
		v.assert.fail(fmt.Sprintf("expected panic with '%v' but got '%v'",
			expected, e))
	}
}

// Catch calls the given function, catching and returning panics
func Catch(f func()) (result any) {
	defer func() {
		if e := recover(); e != nil {
			//dbg.PrintStack()
			result = e
		}
	}()
	f()
	return
}

//-------------------------------------------------------------------

func (a assert) fail(args ...any) {
	if a.t != nil {
		a.t.Helper()
	}
	if len(a.msg) > 0 {
		args = append(append(args, "msg: "), a.msg...)
	}
	s := fmt.Sprintln(args...)
	log.Println("ASSERT FAILED:", s)
	dbg.PrintStack()
	if a.t != nil {
		a.t.Error("\n" + s)
	} else {
		panic("assert failed: " + s)
	}
}

func ShouldNotReachHere() int {
	// return type is so it can be called like panic(ShouldNotReachHere())
	// so the compiler knows it doesn't continue
	e := "ASSERT FAILED: should not reach here"
	log.Println(e)
	dbg.PrintStack()
	panic(e)
}
