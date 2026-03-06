// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("basic", `
123, 123
"return 123", 123
"123 + 456", 579
'"hello" $ " " $ "world"', '"hello world"'
"1 + 2 * 3", "7"
"a = 2; -a", "-2"
"a = true; not a", false
"x = 123; y = 456; x + y", 579

"1000 is 1e3"
"1e3 is 1000"
"10e3 is 1e4"
".001 is 0.001"

"a = 0; ++a", 1
"a = 0; a++", 0
"a = 0; a++; a", 1
"a = 0; --a", -1
"a = 0; a--", 0
"a = 0; a--; a", -1

"Suneido.a = 0; ++Suneido.a", 1
"Suneido.a = 0; Suneido.a++", 0
"Suneido.a = 0; Suneido.a++; Suneido.a", 1
"Suneido.a = 0; --Suneido.a", -1
"Suneido.a = 0; Suneido.a--", 0
"Suneido.a = 0; Suneido.a--; Suneido.a", -1

"Suneido[0] = 0; ++Suneido[0]", 1
"Suneido[0] = 0; Suneido[0]++", 0
"Suneido[0] = 0; Suneido[0]++; Suneido[0]", 1
"Suneido[0] = 0; --Suneido[0]", -1
"Suneido[0] = 0; Suneido[0]--", "0"
"Suneido[0] = 0; Suneido[0]--; Suneido[0]", -1

// integer
"2 < 3", true
"3 < 2", false
"2 < 2", false

// dnum
"1.1 < 1.2", true
"1.2 < 1.1", false

// integer & dnum
".5 < 99", true
"99 < .5", false

// object
"#(.5) < #(99)", true
"#(99) < #(.5)", false

"2 <= 3", true
"3 <= 2", false
"2 <= 2", true

"2 > 3", false
"3 > 2", true
"2 > 2", false

"2 >= 3", false
"3 >= 2", true
"2 >= 2", true

// string
"'' < 'a'", true
"'a' < 'b'", true
"'ab' < 'abc'", true
"'aba' < 'abc'", true

"'abc' =~ '^[a-c]+$'", true
"'abcx' !~ '^[a-c]+$'", true

"s = 1; s $= 2", "'12'"
"n = 10; n -= 5", 5

"t = true; f = false; t and t", true
"t = true; f = false; t and f", false
"t = true; f = false; f and t", false
"t = true; f = false; f and f", false
"t = true; f = false; t or t", true
"t = true; f = false; t or f", true
"t = true; f = false; f or t", true
"t = true; f = false; f or f", false

"x = true; x ? 1 : 2", 1
"x = false; x ? 1 : 2", 2

"0 in (1,2,3)", false
"1 in (1,2,3)", true
"2 in (1,2,3)", true
"3 in (1,2,3)", true
"0 not in (1,2,3)", true
"1 not in (1,2,3)", false
"2 not in (1,2,3)", false
"3 not in (1,2,3)", false

"i = 0; while i < 4 { ++i } i", 4

"x = 4; if x > 3 { return 1 } else { return -1 }", 1
"x = 2; if x > 3 { return 1 } else { return -1 }", -1

"x = 2
switch {
case x < 3: return -1
case x is 3: return 0
case x > 3: return +1
}", -1

"x = 3
switch {
case x < 3: return -1
case x is 3: return 0
case x > 3: return +1
}", 0

"x = 4
switch {
case x < 3: return -1
case x is 3: return 0
case x > 3: return +1
}", 1

"x = 4
switch x {
case 3: return 0
default: return 1
}", 1

"s='hello'; s[0]", 'h'
"s='hello'; s[4]", 'o'
"s='hello'; s[9]", "''"
"s='hello'; s[-1]", 'o' // end relative
"s='hello'; s[-5]", 'h' // end relative
"s='hello'; s[-6]", "''"

"x=#(1,2,3); x[0]", 1
"x=#(1,2,3); x[2]", 3
"x=#(1,2,3); x[3]" throws "member not found: 3"
"x=#(1,2,3); x[-1]" throws "member not found: -1"

// int vs dnum index
"Suneido[1]=123; Suneido[.5 + .5]", 123
"Suneido[1.5 + .5]=456; Suneido[2]", 456

// only "" and false should convert to number or integer

// convert to number - folding
"123 + ''", 123
"123 + false", 123
"-true" throws "can't convert true to number"
"123 + true" throws "can't convert true to number"
"123 + '111'" throws "can't convert String to number"
// convert to number
"x = 123; x + ''", 123
"x = 123; x + false", 123
"x = true; -x" throws "can't convert true to number"
"x = 123; x + true" throws "can't convert true to number"
"x = 123; x + '111'" throws "can't convert String to number"
// convert to integer - folding
"0xff & 0xf", 0xf
"~true" throws "can't convert true to integer"
"0xff & true" throws "can't convert true to integer"
"0xff & '1'" throws "can't convert String to integer"
"4.8 % 2" throws "can't convert number to integer"
// convert to integer
"x = 0xff; x & 0xf", 0xf
"x = true; ~x" throws "can't convert true to integer"
"x = 0xff; x & true" throws "can't convert true to integer"
"x = 0xff; x & '1'" throws "can't convert String to integer"
"x = 4.8; x % 2" throws "can't convert number to integer"

"s = 'hello'; s[1.2]" throws "member not found"
"s = 'hello'; s['']" throws "member not found"
"s = 'hello'; s[false]" throws "member not found"
"s = 'hello'; s['' ..]" throws "indexes must be integers"
"s = 'hello'; s[false ..]" throws "indexes must be integers"
"s = 'hello'; s['' ::]" throws "indexes must be integers"
"s = 'hello'; s[false ::]" throws "indexes must be integers"
"x = #(); x['' ::]" throws "indexes must be integers"
"x = #(); x[false ::]" throws "indexes must be integers"

// but second value in ranges should accept "" and false
"s = 'hello'; s[..'']", "''"
"s = 'hello'; s[..false]", "''"
"s = 'hello'; s[::'']", "''"
"s = 'hello'; s[::false]", "''"

"x = Object(1,2,3); x[0::2]", "#(1,2)"
"x = Object(3); x[0::1][0] = 9; x", "#(3)"

// folding/optimization still type checks

'0 + "123"' throws "can't convert"
'"123" + 0' throws "can't convert"
'x = "123"; 0 + x' throws "can't convert"
'x = "123"; x + 0' throws "can't convert"

'1 * "123"' throws "can't convert"
'"123" * 1' throws "can't convert"
'x = "123"; 1 * x' throws "can't convert"
'x = "123"; x * 1' throws "can't convert"

'"" $ #()' throws "can't convert"
'#() $ ""' throws "can't convert"
'x = #(); x $ ""' throws "can't convert"
'x = #(); "" $ x' throws "can't convert"
`)

var _ = Register("argument handling", `
"Object().Add(@Seq(1000)).Size()", 1000 //  // Add(@ shouldn't expand onto stack

"Object(a: 1, a: 2)" throws "duplicate argument name"

"a=1; Object(:a)", "#(a: 1)"

"a=1; Object([:a])", "#((a: 1))"
"b=1; Object(a: [:b])", "#(a: [b: 1])"
"Object(@#(1, 2, a: 3, b: 4))", "#(1, 2, a: 3, b: 4)"
"Object(@+1#(1, 2, a: 3, b: 4))", "#(2, a: 3, b: 4)"
"x=[];y=Object(@x);x.a=1;y", "#()"
"Record(@#(1, 2, a: 3, b: 4))", "#{1, 2, a: 3, b: 4}"
"Record(@+1#(1, 2, a: 3, b: 4))", "#{2, a: 3, b: 4}"
"x=[];y=Record(@x);x.a=1;y", "#{}"

'function(x){x}(@#(123))', 123
'function(x){x}(@+1#(123,456))', 456
'function(x){x}(@[123])', 123
`)

var _ = Register("unary plus and minus conversions", `
"x = true; +x" throws "can't convert true to number"
"x = false; +x", 0
"x = true; -x" throws "can't convert true to number"
"x = false; -x", 0
// fold
"+true" throws "can't convert true to number"
"+false", 0
"-true" throws "can't convert true to number"
"-false", 0
`)

var _ = Register("calls", `
// different numbers of arguments
"function(){123}()", 123
"function(a){a}(123)", 123
"function(a,b){a+b}(1,2)", 3
"function(a,b,c){a+b+c}(1,2,3)", 6
"function(a,b,c,d){a+b+c+d}(1,2,3,4)", 10
"function(a,b,c,d,e){a+b+c+d+e}(1,2,3,4,5)", 15

// no params
"function(){123}()", 123
"function(){123}(a: 1)", 123
"function(){123}(a: 1, b: 2)", 123
"function(){123}(@#())", 123
"function(){123}(@#(a: 1, b: 2))", 123
"function(){123}(@+1#())", 123
"function(){123}(@+1#(a: 1, b: 2))", 123
"function(){123}(1)" throws "too many arguments"

// @param
"function(@x){x}(123)", "#(123)"
"function(@x){x}(1,2,3)", "#(1, 2, 3)"
"function(@x){x}(1,2,a:3,b:4)", "#(1, 2, a: 3, b: 4)"

// some params
"function(a,b){a+b}()" throws "missing argument"
"function(a,b){a+b}(1)" throws "missing argument"
"function(a,b){a+b}(1,2)", 3
"function(a,b){a+b}(@#(1,2))", 3
"function(a,b){a+b}(@+1#(1,2,3))", 5
"function(a,b){a+b}(@#(1,2,x:3,y:4))", 3
"function(a,b){a+b}(1,2,3)" throws "too many arguments"
"function(a,b){a+b}(@#(1,2,3))" throws "too many arguments"

// dynamic param defaults
"_a=1; _b=2; function(_a,_b){[a,b]}()", "[1,2]"
"_a=1; _b=2; function(_a,_b){[a,b]}(3)", "[3,2]"
"_a=1; _b=2; function(_a,_b){[a,b]}(3,4)", "[3,4]"

// non-method member call
"ob = Object(f: function(){123}); ob.f()" throws "method not found"
"ob = Object(f: function(){123}); (ob.f)()", 123

// return nil passthrough
'f = function(){}; return f()', nil
'f = function(){}; x = true; return x ? f() : f()', nil
'f = function(){}; x = false; return x ? f() : f()', nil
'f = function(){}; b = { return f() }; b()', nil // closure
'b = { f = function(){}; return f() }; b()', nil // non-closure
`)

var _ = Register("naming", `
// test Def compiles string
"Def(#Foo, 'function(){}')
Type(Foo)", "'Function'"

"foo = function(){}
Name(foo)", "foo"

"c = class { Foo(){} }
Name(c.Foo).Has?('.Foo')"

"Def(#Foo, 'function(){}')
Name(Foo)", "Foo"

"Def(#Foo, 'class{}')
Name(Foo)", "Foo"
`)

var _ = Register("blocks", `
// non closures

'b = { 123 }; b()', 123

'b = {|x,y| x+y }; b(1,2)', 3

'b = {|it| 2 * it }; b(3)', 6

'b = { 2 * it }; b(3)', 6

// closures

'x = 6; y = 3; b = { x / y }; b()', 2

'x = 6; b = {|a| x / a }; b(3)', 2

'x = 123; b = { x = 456 }; b(); x', 456

'class { x: 123; F() { b = { .x }; b() } }.F()', 123 // this

// escaping closures

'f = function (x) { return {|a| x * a } }; b = f(2); b(3)', 6

'f = function (x,ob) { ob.b = {|a| x * a };0 };
ob = Object(); f(2,ob); (ob.b)(3)', 6

// parameters

't = false; b = {|t| return t }; b(true)', 'true'

// block parameter names offset from parent function names
'outside = 0
b = {|env| outside }
env = 123
b(env: env)', 0

// nested blocks (check that locals are shared)
'run = function (block) { block() }
run() { run() { x = 123 } }
x', 123
`)

var _ = Register("for-in", `
'ob = Object()
for (iter = #(2,4,8).Iter(); iter isnt x = iter.Next(); )
	ob.Add(x)
ob', '#(2,4,8)'

'ob = Object()
for x in #(2,4,8)
	ob.Add(x)
ob', '#(2,4,8)'

'ob = Object()
for x in #{2,4,8}
	ob.Add(x)
ob', '#(2,4,8)'
`)

var _ = Register("try-catch-throw", `
'x = 0; try x = 123; x', 123

'x = 0; try x = 123 catch x = 456; x', 123

'try return 123', 123

'throw "foo"' throws "foo"

'try throw "foo"; 123', 123

'try throw "foo" catch {}; 123', 123

'try throw "foo" catch (e) return "e=" $ e', '"e=foo"'

'try throw "foo" catch (e, "x") return "e=" $ e' throws "foo"

'try throw "foo" catch (e, "f") return "e=" $ e', '"e=foo"'

'try throw "foobar" catch (e, "*bar") return "e=" $ e', '"e=foobar"'

'try throw "foo" catch (e, "f") {}; try throw "foobar" catch (e, "x") {}',
	throws "foobar"

'try throw "foo" catch (e) return Type("e=" $ e)', 'Except'

'try throw "foo" catch (e) return Type(e $ "=e")', 'Except'

'try throw "foo" catch (e) return Type("x".Repeat(1000) $ "y" $ e)', 'Except'

'try throw "error" catch (e) return "error" is e', true

'try x=123
catch (e) throw "caught"
throw "after"' throws "after"

// return from block returns from parent
'b = { return 123 }; b(); 456', 123

// including indirectly (c makes f also a catching block parent)
'f = function (b){ c={}; b() }; f({ return 123 }); 456', 123

'Finally({ throw "foo" }, { throw "bar" })' throws 'foo'

'Finally({ 123 }, { throw "bar" })' throws 'bar'
`)

var _ = Register("sequences", `
'#(2, 4, a: 6).Values()', '#(2, 4, 6)'

'#(2, 4, a: 6).Members()', '#(0, 1, a)'

'#(2, 4, a: 6).Assocs()', '#((0,2), (1,4), (a,6))'

'Seq(5)', '#(0,1,2,3,4)'

'Seq(1, to: 5)', '#(1,2,3,4)'

'Seq(1, to: 10, by: 2)', '#(1,3,5,7,9)'

'Def("Q", "class {
	Next() { .i < 4 ? .i++ : this }
	Dup() { new (.Base()) }
	Infinite?() { false }
	i: 0 }"); true'

'Sequence(Q())', '#(0,1,2,3)'

'q = Sequence(Q()); q.Infinite?()', false

'q = Sequence(Q()); q.Instantiated?()', false

'q = Sequence(Q()); q.Infinite?(); q.Instantiated?()', false

'q = Sequence(Q()); q.Iter().Next()', 0

'q = Sequence(Q()); q.Iter().Next(); q.Instantiated?()', false

// two dups = instantiate
'q = Sequence(Q());
q.Iter().Next(); q.Iter().Next(); q.Instantiated?()', true

'q = Sequence(Q()); q.Join(",")', '"0,1,2,3"'

'q = Sequence(Q()); q.Join(","); q.Instantiated?()', false

'ob = Object()
for x in Sequence(Q())
	ob.Add(x)
ob', '#(0,1,2,3)'

// infinite sequence
'Def("Qi", "class {
	Next() { .i++ }
	Dup() { new (.Base()) }
	Infinite?() { true }
	i: 0 }"); true'

'q = Sequence(Qi()); q.Infinite?()', true

'q = Sequence(Qi()); q.Instantiated?()', false

'q = Sequence(Qi()); q.Infinite?(); q.Instantiated?()', false

// infinite, so two dups != instantiate
'q = Sequence(Qi());
q.Iter().Next(); q.Iter().Next(); q.Instantiated?()', false

'Seq?(123)', false

'Seq?(#())', false

'Seq?(Seq())', true

'Seq?(Sequence(Q()))', true
`)

var _ = Register("Type", `
"Type(true)", "'Boolean'"
"Type('foo')", "'String'"
"Type(123)", "'Number'"
"Type(function () { })", "'Function'"
"Type(class { })", "'Class'"
"Type(class { }())", "'Instance'"
"Type(Pack)", "'BuiltinFunction'"
"Type(Date)", "'BuiltinFunction'"
"Type(Date())", "'Date'"
"Type(Adler32)", "'BuiltinFunction'"
"Type(Adler32())", "'BuiltinClass'"
"Type(Md5)", "'BuiltinFunction'"
"Type(Md5())", "'BuiltinClass'"
"Type(Seq)", "'BuiltinFunction'"
"Type(Seq())", "'Object'"
`)

var _ = Register("Display", `
"Display(Pack)", "'Pack /* builtin function */'"
"Display(Seq)", "'Seq /* builtin function */'"
"Display(Seq())", "'infiniteSequence'"
"Display(Date)", "'Date /* builtin class */'"
"Display(Thread)", "'Thread /* builtin class */'"
"Display(Adler32)", "'Adler32 /* builtin function */'"
"Display(Md5)", "'Md5 /* builtin function */'"
"Display(Adler32())", "'adler32'"
"Display(Md5())", "'md5'"
`)
