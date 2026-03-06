// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("strings string.Split", `
'"hello".Split()' '#(h, e, l, l, o)'
'"one,two,three".Split(",")', '#(one,two,three)'
'"".Split(",")', '#()'
'"one".Split(",")', '#(one)'
'"one,".Split(",")', '#(one)'
'",one".Split(",")', '#("",one)'
'"one=>two=>three".Split("=>")', '#(one,two,three)'
`)

var _ = Register("strings string.FindLast", `
// without position
'"this is a test".FindLast("")', 14
'"this is a test".FindLast("this")', 0
'"this is a test".FindLast("test")', 10
'"this is a test".FindLast("is")', 5
'"this is a test".FindLast("xyz")', false

// with position
'"this is a test".FindLast("", 0)', 0
'"this is a test".FindLast("", 5)', 5
'"this is a test".FindLast("", 14)', 14
'"this is a test".FindLast("", 99)', 14
'"this is a test".FindLast("this", 9)', 0
'"this is a test".FindLast("this", 4)', 0
'"this is a test".FindLast("this", 99)', 0
'"this is a test".FindLast("this", 2)', 0
'"this is a test".FindLast("this", 0)', 0
'"this is a test".FindLast("this", -1)', false
'"this is a test".FindLast("this", -9)', false
'"this is a test".FindLast("test", 14)', 10
'"this is a test".FindLast("is", 9)', 5
'"this is a test".FindLast("is", 5)', 5
'"this is a test".FindLast("is", 4)', 2
'"this is a test".FindLast("xyz", 99)', false
'"this is a test".FindLast("xyz", -9)', false
'"this is a test".FindLast("xyz", 0)', false
`)

var _ = Register("strings string.Entab/Detab", `
'"".Entab()', "''"
'"foo".Entab()', "'foo'"
'"foo bar".Entab()', "'foo bar'"
'"    foo".Entab()', "'\tfoo'"
'"  \tfoo".Entab()', "'\tfoo'"
'" \t foo".Entab()', "'\t foo'"
'" \t foo  \t  ".Entab()', "'\t foo'"
'"foo\tbar".Entab()', "'foo\tbar'"
'"    foo\r\n    bar\r\n".Entab()', "'\tfoo\r\n\tbar\r\n'"
'"foo\t\nbar".Entab()', "'foo\nbar'"
'"foo\t\r\nbar".Entab()', "'foo\r\nbar'"
'"\t".Entab()', "''"
'"\t\n".Entab()', "'\n'"
'"\t\r\n".Entab()', "'\r\n'"

'"".Detab()', "''"
'"foo bar".Detab()', "'foo bar'"
'"  foo".Detab()', "'  foo'"
'"\tfoo".Detab()', "'    foo'"
'"  \tfoo".Detab()', "'    foo'"
'"\t\tfoo".Detab()', "'        foo'"
'" \t \tfoo".Detab()', "'        foo'"
'"x\ty".Detab()', "'x   y'"
'"\tfoo\n\tbar".Detab()', "'    foo\n    bar'"
'"\tfoo\r\n\tbar".Detab()', "'    foo\r\n    bar'"
`)

var _ = Register("strings string.Match", `
'"now is the time".Match("foo")'			false
'"now is the time".Match("now")'			"#(#(0,3))"
'"now is the time".Match("the")'			"#(#(7,3))"
'"now is the time".Match("time")'			"#(#(11,4))"

'"now is the time".Match("foo", prev:)'		false
'"now is the time".Match("now", prev:)'		"#(#(0,3))"
'"now is the time".Match("the", prev:)'		"#(#(7,3))"
'"now is the time".Match("time", prev:)'	"#(#(11,4))"

'"big bigger biggest".Match("big")'			"#(#(0,3))"
'"big bigger biggest".Match("big", prev:)'	"#(#(11,3))"

'"big bigger biggest".Match("big", 1)'			"#(#(4,3))"
'"big bigger biggest".Match("big", 10, prev:)'	"#(#(4,3))"

'"now is the time".Match("(\w+) (\w+) (\w+)")'	"#((0,10),(0,3),(4,2),(7,3))"
`)

var _ = Register("strings string.Extract", `
'"hello world".Extract(".....$")', '"world"'
'"hello world".Extract("w(..)ld")', '"or"'
'"hello world".Extract("(\w+) (\w+)")', '"hello"'
'"hello world".Extract("(\w+) (\w+)", 0)', '"hello world"'
'"hello world".Extract("(\w+) (\w+)", 1)', '"hello"'
'"hello world".Extract("(\w+) (\w+)", 2)', '"world"'
`)

var _ = Register("strings string.Eval", `
'"123".Eval()', 123
'Def("X", 123); "X".Eval()', 123
'"123 + 456".Eval()', 579
'"return".Eval()', '""'
'"123 + 456".Eval2()', '#(579)'
'"return".Eval2()', '#()'
`)

var _ = Register("strings string.Asc/Chr", `
'"\x00".Asc()', '0'
'"a".Asc()', '97'
'"\xff".Asc()', '255'
'0.Chr()', '"\x00"'
'97.Chr()', '"a"'
'255.Chr()', '"\xff"'
`)

var _ = Register("strings calling a string", `
'#Upper("abc")', '"ABC"'
'#Upper(@#(abc))', '"ABC"'
'#Upper(@+1#(abc,def))', '"DEF"'
`)
