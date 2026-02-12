package tests

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestReturnConsistency(t *testing.T) {
	test := func(src string, expected ...string) {
		t.Helper()
		// fmt.Println("----------------")
		// fmt.Println(src)
		// if expected != nil {
		// 	fmt.Println(expected)
		// }
		_, results := compile.Checked(nil, src)
		assert.This(results).Is(expected)
	}

	// Consistent returns - should pass
	test("function () { return }")
	test("function () { return 1 }")
	test("function () { return 1, 2 }")
	test("function (a) { if a \n return 1 \n else \n return 2 }")
	test("function (a) { if a \n return \n else \n return }")
	test("function (a) { if a \n return 1, 2 \n else \n return 3, 4 }")

	// unknown number of return values from function call
	test("function (a) { if a \n return \n return Date() }")

	// Inconsistent explicit returns - should fail
	test("function (a) { if a \n return \n else \n return 1 }",
		"WARNING: inconsistent number of return values @38")
	test("function (a) { if a \n return 1 \n else \n return 1, 2 }",
		"WARNING: inconsistent number of return values @40")
	test("function (a) { if a \n return 1, 2 \n else \n return }",
		"WARNING: inconsistent number of return values @43")

	// Consistent implicit returns - should pass
	test("function () { 123 }")          // implicit return 1 value
	test("function (a) { if a \n a() }") // implicit return 0 values
	test("function () { }")              // implicit return 0 values
	test("function (a) { if a \n return \n else \n a() }")

	// Inconsistent implicit vs explicit returns - should fail
	test("function (a) { if a \n return 1 }",
		"WARNING: inconsistent number of return values @30")
	test("function (a) { if a \n return 1, 2 \n else \n a = 1 }",
		"WARNING: inconsistent number of return values @48")
	test("function (a) { if a \n a = 1 \n else \n return }")

	// Multiple inconsistent returns - should fail on first
	test("function (a) { if a \n return 1 \n else if a \n return 1, 2 \n else \n return }",
		"WARNING: inconsistent number of return values @45",
		"WARNING: inconsistent number of return values @66")

	// Function call returns - should be allowed (no consistency checking)
	test("function (f) { return f() }")                                     // function call can return any number
	test("function (a,f) { if a \n return f() \n else \n return }")         // mixed with regular return - should pass
	test("function (a,f) { if a \n return \n else \n return f() }")         // mixed with regular return - should pass
	test("function (a,f,g) { if a \n return f() \n else \n return g() }")   // multiple function calls - should pass
	test("function (a,f) { if a \n return f() \n else \n return 1, 2, 3 }") // function call mixed with multiple values - should pass

	// Block returns (return from blocks returns from parent function)
	test("function (f) { f({ return 1 }) }")                   // block return should share parent's count
	test("function (f) { if f \n return \n f({ return 2 }) }", // block return inconsistent with parent
		"WARNING: inconsistent number of return values @35")
	test("function (f) { f({ return }) \n return 1 }", // block return inconsistent with parent
		"WARNING: inconsistent number of return values @31")
	test("function (f) { f({ return 1 }) \n return 2 }") // block return consistent with parent

	// Block implicit returns should not be checked
	test("function (f) { f({ 123 }) }")             // block implicit return - no checking
	test("function (f) { f({ 123 }) \n return 1 }") // block implicit return with parent explicit return
}
