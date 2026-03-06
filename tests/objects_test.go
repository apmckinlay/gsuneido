// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("objects", `
'Type([])', Record
'Type([a: 1, b: 2])', Record
'Type([1, 2, a: 3])', Object
'x=1; Type([x, a: 3])', Object

'Object?(Object())'
'Object?(#())'
'Object?(Record())'
'Object?(#{})'
'Object?([])'

'x = #()
x.a = 123' throws "readonly"

'x = #{}
x.a = 123' throws "readonly"

'x = Object()
x.Set_readonly()
x.a = 123' throws "readonly"

'x = Record()
x.Set_readonly()
x.a = 123' throws "readonly"

// range
'#(1,2,3)[1..]', '#(2,3)'
'#(1,2,3, a: 4, b: 5)[1..]', '#(2,3)'
'#(1,2,3)[1::]', '#(2,3)'
'#(1,2,3, a: 4, b: 5)[1::]', '#(2,3)'

// recursive
'x = Object(y: Object())
x.Set_readonly()
x.y.a = 123' throws "readonly"

'x = Record()
x.a', "''"

'x = Object()
x.a' throws "member not found"

'x = Object()
x.Set_default(0)
x.a', 0

'x = Record()
x.Set_default(0)
x.a', 0

'x = #()
x.Set_default(0)' throws "readonly"

'x = #{}
x.Set_default(0)' throws "readonly"

'x = Object()
x.Set_readonly()
x.Set_default(0)' throws "readonly"

'x = Object()
x.Set_default(#())
d = x.foo
d.m = 123
x.foo', '#(m: 123)'

'x = Object()
x.Set_default(#{})
Type(x.foo)', 'Record'

'x = Object()
x.Set_default(#())
d = x.foo
d.m = 123
x.bar', '#()'

'x = Record()
x.Set_default(#())
d = x.foo
d.m = 123
x.foo', '#(m: 123)'

'x = Record()
x.Set_default(#())
d = x.foo
d.m = 123
x.bar', '#()'

// object.Size

'#(1,2,3,a:4,b:5).Size()', 5

'#(1,2,3,a:4,b:5).Size(list:)', 3

'#(1,2,3,a:4,b:5).Size(named:)', 2

'#(1,2,3,a:4,b:5).Size(list:, named:)', 5

'#(1,2,3).Size(named:)', 0

'#(a:4,b:5).Size(list:)', 0

// object.GetDefault

'#(a:123,b:456).GetDefault(#a, 0)', 123

'#(a:123,b:456).GetDefault(#c, 789)', 789

'#(a:123,b:456).GetDefault(#c, { 123 + 456 })', 579

'#(a:123,b:456).GetDefault(#c){ 123 + 456 }', 579

// object.Add

'Object().Add(1)', '#(1)'

'Object(1, 2).Add(3, 4)', '#(1, 2, 3, 4)'

'Object().Add(at: "a")', '#()'

'Object().Add(1 at: "a")', '#(a: 1)'

'Object().Add(1, 2 at: "a")' throws "position"

'Object().Add(1 at: "a", x: 2, y: 3)', '#(a: 1)' // ignore extra named

'Object().Add(1 at: 9)', '#(9: 1)'

'Object(1, 4).Add(2, 3 at: 1)', '#(1, 2, 3, 4)'

// same calls with @args

'Object().Add(@#(1))', '#(1)'

'Object(1, 2).Add(@#(3, 4))', '#(1, 2, 3, 4)'

'Object().Add(@#(at: "a"))', '#()'

'Object().Add(@#(1 at: "a"))', '#(a: 1)'

'Object().Add(@#(1, 2 at: "a"))' throws "position"

'Object().Add(@#(1 at: "a", x: 2, y: 3))', '#(a: 1)' // ignore extra named

'Object().Add(@#(1 at: 9))', '#(9: 1)'

'Object(1, 4).Add(@#(2, 3 at: 1))', '#(1, 2, 3, 4)'

// same calls with @+1 args

'Object().Add(@+1 #(0, 1))', '#(1)'

'Object(1, 2).Add(@+1 #(0, 3, 4))', '#(1, 2, 3, 4)'

'Object().Add(@+1 #(1 at: "a"))', '#()'

'Object().Add(@+1 #(0, 1 at: "a"))', '#(a: 1)'

'Object().Add(@+1 #(0, 1, 2 at: "a"))' throws "position"

'Object().Add(@+1 #(0, 1 at: "a", x: 2, y: 3))', '#(a: 1)' // ignore extra named

'Object().Add(@+1 #(0, 1 at: 9))', '#(9: 1)'

'Object(1, 4).Add(@+1 #(0, 2, 3 at: 1))', '#(1, 2, 3, 4)'

// Sort

'#().Sort!()' throws "readonly"

'Object().Sort!()', '#()'

'Object(1,2,3,4).Sort!()', '#(1,2,3,4)'

'Object(2,4,1,3).Sort!()', '#(1,2,3,4)'

'Object(2,4,1,3).Sort!(function (x,y) { x > y })', '#(4,3,2,1)'

'Record(2,4,1,3).Sort!(function (x,y) { x > y })', '#{4,3,2,1}'

// Delete

'[11,22,a:33,b:44].Delete(#foo)', '[11,22,a:33,b:44]'

'[11,22,a:33,b:44].Delete(all:)', '[]'

'[11,22,a:33,b:44].Delete(0)', '[22,a:33,b:44]'

'[11,22,a:33,b:44].Delete(1)', '[11,a:33,b:44]'

'[11,22,a:33,b:44].Delete("a")', '[11,22,b:44]'

'[11,22,a:33,b:44].Delete(#b)', '[11,22,a:33]'

'[11,22,a:33,b:44].Delete(1,0,#a,"b")', '[]'

// Erase

'[11,22,a:33,b:44].Erase(#foo)', '[11,22,a:33,b:44]'

'[11,22,a:33,b:44].Erase(0)', '[1: 22,a:33,b:44]'

'[11,22,a:33,b:44].Erase(1)', '[11,a:33,b:44]'

'[11,22,a:33,b:44].Erase("a")', '[11,22,b:44]'

'[11,22,a:33,b:44].Erase(#b)', '[11,22,a:33]'

'[11,22,a:33,b:44].Erase(1,0,#a,"b")', '[]'

// Join

'[11,22,33,foo:44].Join()', '"112233"'

'[11,22,33,foo:44].Join("=>")', '"11=>22=>33"'

// Find

'[11,22,a:33,b:44].Find(99)', false

'[11,22,a:33,b:44].Find(0)', false

'[11,22,a:33,b:44].Find(#a)', false

'[11,22,a:33,b:44].Find(11)', 0

'[11,22,a:33,b:44].Find(22)', 1

'[11,22,a:33,b:44].Find(33)', '#a'

'[11,22,a:33,b:44].Find(44)', '#b'

// Copy

'#().Copy()', '#()'

'#(1,2,3).Copy()', '#(1,2,3)'

'#(11,22,a:33,b:44).Copy()', '#(11,22,a:33,b:44)'

'#(a:33,b:44).Copy()', '#(a:33,b:44)'

'#().Add(123)' throws "readonly"

'#().Copy().Add(123)', '#(123)' // copy of read-only is not read-only

'Type([].Copy())', '"Record"' // copy of a Record is a Record

// Eval

'x = #(A:123).Eval(function(){ })' throws 'no return value'

'#(A:123).Eval(function(){.A})', 123

'#(a:123).Eval(function(a){a + .a}, 456)', 579

'#(a:123).Eval(@#(function(a){a + .a}, 456))', 579

'c = class{A: 222; F(){.A}}; #(A:123).Eval(c.F)', 123 // bound method

// BinarySearch

'#().BinarySearch(123)', 0

'#(1, 2, 3, 3, 4, 5).BinarySearch(0)', 0

'#(1, 2, 3, 3, 4, 5).BinarySearch(1)', 0

'#(1, 2, 3, 3, 4, 5).BinarySearch(2)', 1

'#(1, 2, 3, 3, 4, 5).BinarySearch(3)', 2

'#(1, 2, 3, 3, 4, 5).BinarySearch(4)', 4

'#(1, 2, 3, 3, 4, 5).BinarySearch(9)', 6

'#(1, 2, 3, 3, 4, 5).BinarySearch(2, {|x,y| x <= y})', 2

'#(1, 2, 3, 3, 4, 5).BinarySearch(3, {|x,y| x <= y})', 4

'#(1, 2, 3, 3, 4, 5).BinarySearch(4, {|x,y| x <= y})', 5
`)
