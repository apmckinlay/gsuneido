// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPlay(t *testing.T) {
	pat := Compile(`(\w+) (\w+)`)
	fmt.Println(pat)
	fmt.Println(Match(pat, "Foo bar"))
}

func TestMatch(t *testing.T) {
	yes := func(pat string, str string) {
		assert.T(t).True(Match(Compile(pat), str))
	}
	no := func(pat string, str string) {
        assert.T(t).False(Match(Compile(pat), str))
    }
	yes("a", "a")
	no("a", "")
	no("a", "b")

	yes(".", "a")
	no(".", "")

	yes("a|b", "a")
	yes("a|b", "b")
	no("a|b", "")
	no("a|b", "c")

	yes("a?", "")
	yes("a?", "a")
	no("a?", "b")
}

func ExampleCompile() {
	test := func(rx string) {
		fmt.Printf("/%v/\n%v\n", rx, Compile(rx))
	}
	test("abc")
	test("a|b")
	test("ab?c")
	test("ab+c")
	test("ab*c")
	// Output:
	// /abc/
	// 0: Save 0
	// 2: Char a
	// 4: Char b
	// 6: Char c
	// 8: Save 1
	// 10: Stop
	//
	// /a|b/
	// 0: Save 0
	// 2: SplitFirst 10
	// 5: Char a
	// 7: Jump 12
	// 10: Char b
	// 12: Save 1
	// 14: Stop
	//
	// /ab?c/
	// 0: Save 0
	// 2: Char a
	// 4: SplitFirst 9
	// 7: Char b
	// 9: Char c
	// 11: Save 1
	// 13: Stop
	//
	// /ab+c/
	// 0: Save 0
	// 2: Char a
	// 4: Char b
	// 6: SplitFirst 4
	// 9: Char c
	// 11: Save 1
	// 13: Stop
	//
	// /ab*c/
	// 0: Save 0
	// 2: Char a
	// 4: SplitFirst 12
	// 7: Char b
	// 9: Jump 7
	// 12: Char c
	// 14: Save 1
	// 16: Stop
}
