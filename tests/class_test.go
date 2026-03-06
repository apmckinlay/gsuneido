// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("classes and instances", `
// CallClass
"class{CallClass(){123}}()", 123
"class{CallClass(a){a}}(123)", 123
"class{CallClass(a,b){a+b}}(1,2)", 3
"class{CallClass(a,b,c){a+b+c}}(1,2,3)", 6
"class{CallClass(a,b,c,d){a+b+c+d}}(1,2,3,4)", 10
"class{CallClass(a,b,c,d,e){a+b+c+d+e}}(1,2,3,4,5)", 15

// class "static" method
"class{M(){123}}.M()", 123
"class{M(a){a}}.M(123)", 123
"class{M(a,b){a+b}}.M(1,2)", 3
"class{M(a,b,c){a+b+c}}.M(1,2,3)", 6
"class{M(a,b,c,d){a+b+c+d}}.M(1,2,3,4)", 10
"class{M(a,b,c,d,e){a+b+c+d+e}}.M(1,2,3,4,5)", 15

// access to "this"
"class{X: 123}.X", 123
"class{M(){.X} X: 123}.M()", 123
"class{M(){.x} x: 123}.M()", 123
"class{CallClass(){.x} x: 123}()", 123

'// class without base
class { }'

'// class with base
c = class : Test { }
c = Test { }'

'// comment
Test
	{
	}'

'// Type of class
c = class { }
Type(c) is "Class"'

'// class equality is identity
c = class { }
c is c and c isnt class { }'

'// class get
c = class { M: 123 }
c.M is 123'

'// class get inherits
Def(#A, class { X: 123 })
c = A { }
c.X is 123'

'// classes are readonly()
c = class { }
c.Mem = 123',
	throws ""

'// class members()
c = class { M: 123, F() { } }
c.Members().Sort!() is #(F, M)'

'// class methods()
c = class { F() { 123 } }
c.F() is 123'

'// private member
c = class { x: 123 }
c.x' throws "member not found"

'// private member accessible to method
c = class { x: 123; F() { .x } }
c.F() is 123'

'class { F() { this = 5 }}' throws "error"

'class { F() { super = 5 }}' throws "error"

'// instance
c = class { }
new c'

'// cannot create instance of instance
x = class{}()
new x' throws "can't create instance of instance"

'// Type of instance
c = class { }
i = new c
Type(i) is "Instance"'

'// instances inherit members
c = class { M: 123 }
i = new c
i.M is 123'

'// memberq includes inherited()
c = class { M: 123 }
i = new c
i.Member?("M")'

'// method lookup starts in class()
c = class { F() { 123 } }
i = new c
i.F = "foo"
c.F() is 123'

'// classes inherit methods
Def(#A, class { F() { 123} })
Def(#B, A { })
c = B{}
c.F() is 123'

'// instances inherit methods
Def(#A, class { F() { 123} })
Def(#B, A { })
c = B{}
x = c()
x.F() is 123'

'// instance equality is value()
c = class { UseDeepEquals: true }
new c is new c'

'// default CallClass is new()
c = class { UseDeepEquals: true }
c() is new c'

'// instances are modifiable()
c = class { }
i = new c
i.Mem = 123
i.Mem is 123'

'// New
c = class { New(x) { .X = x } }
c(123).X is 123'

'// New is chained
Def(#A, class { New() { .A = 1 }})
Def(#B, A { New() { .B = 2 }})
x = new B{ New() { .C = 3 }}
x.A is 1 and x.B is 2 and x.C is 3'

'// super requires parent
class { F() { super.F() }}' throws "super requires parent"

'// super requires parent
class { New() { super() }}' throws "super requires parent"

'// super(...) must be first
Def(#A, class { F(x) { x + 1 }})
A { New() { Other(); super() }}' throws "super"

'// super(...) only valid in New
Def(#A, class { F(x) { x + 1 }})
x = new A { F() { super(1) }}' throws "super"

'// super New call
Def(#A, class { New(x){ .X = x }})
x = new A { New() { super(123) }}
x.X is 123'

'// super method call
Def(#A, class { F(x) { x + 1 }})
x = new A { F() { super.F(1) }}
x.F() is 2'

'// New with private dot param
c = class { New(.x) { } F() { .x }}
c(123).F() is 123'

'// New with public dot param
c = class { New(.X) { } }
c(123).X is 123'
`)

var _ = Register("privatization", `
"Def(#Foo, 'class { x: 123 }')
Foo.Foo_x", 123

"Def(#Foo, 'class { f(){ 123 } }')
Foo.Foo_f()", 123

"Def(#Foo, 'class { F(x){ .x = x } }')
ob = Foo()
ob.F(123)
ob.Foo_x", 123

"Def(#Foo, 'class { F(.x){ } }')
ob = Foo()
ob.F(123)
ob.Foo_x", 123

`)

var _ = Register("class getter methods", `
'// normal member
c = class { A: 123 }
c.A is 123 and c().A is 123'

'// public getter
c = class { Getter_A() { 123 } }
c.A is 123 and c().A is 123'

'// private getter
c = class { getter_a() { 123 }; F() { .a } }
c.F() is 123 and c().F() is 123'

'// invalid getter
class { getter_() { } }',
throws "invalid getter"

'// invalid getter
class { getter_A() { } }',
throws "invalid getter"

'// invalid getter
class { Getter_a() { } }',
throws "invalid getter"

'// invalid explicit getter call
c = class { getter_a() { 123 }; F() { .getter_a() }}
c.F()',
throws "method not found"

'// general getter
c = class { Getter_(m) { "(" $ m $ ")" } }
c.Foo is "(Foo)" and c().Foo is "(Foo)"'

'// public getter
i = class { X: 123; Getter_A() { .X } }()
i.X = 456;
i.A', 456

'// public getter
i = class { X: 123; Getter_(m) { .X } }()
i.X = 456;
i.A', 456

'// getter should handle non-string (class)
c = class { Getter_(m) { m } }
c[123]', 123

'// getter should handle non-string (instance)
i = class { Getter_(m) { m } }()
i[123]', 123
`)

var _ = Register("bound methods", `
'c = class { x: 123; F() { .x }}
m = c.F; m()', 123

'c = class { New() { .x = 123 }; F() { .x }}
m = c().F; m()', 123

'c = class { x: 123; f() { .x }; F() { .f }}
m = c.F(); m()', 123

'c = class { New() { .x = 123 }; f() { .x }; F() { .f }}
m = c().F(); m()', 123

'c = class { F(a,b) { } }
m = c.F; m.Params()', '"(a,b)"'

'c = class { F(a,b) { } }
c.F is c.F'
`)

var _ = Register("ToString", `
'Display(class { })', '"/* class */"'

// ToString only affects instances
'Display(class { ToString() { "*foo*" } })', '"/* class */"'

'Def(#C, "class{}"); Display(C())', '"C()"'

'x = class { ToString() { .S }; S: "foo" }()
Display(x)', '"foo"'

'Display(class { ToString() {} }())' throws 'ToString should return a string'

'Display(class { ToString() { #() } }())' throws 'ToString should return a string'
`)

var _ = Register("class/instance methods", `
'c = class {}
c.Size()', 0

'c = class {x:, y:}
c.Size()', 2

'x = class {x:, y:}()
x.Size()', 0

'x = class {x:, y:}()
x.Foo = 123
x.Size()', 1

'x = class{}()
x.Foo = 123
x.Delete(#Foo)
x.Members()', '#()'

'x = class{}()
x.Foo = x.Bar = 123
x.Delete(all:)
x.Members()', '#()'

'class{}.Base()', false

'class{}.Base?(123)', false

'Def("C", "class { }")
C.Base()', false

'Def("C", "class { }")
C.Base?(C)'

'Def("C", "class { }")
C().Base() is C'

'Def("C", "class { }")
C().Base?(C)'

'Def("B", "class { }")
Def("C", "B { }")
C.Base() is B'

'Def("B", "class { }")
Def("C", "B { }")
C().Base() is C'

'Def("B", "class { }")
Def("C", "B { }")
C().Base?(C)'

'Def("B", "class { }")
Def("C", "B { }")
C().Base?(B)'

// Method?

'class{ F(){} }.Method?(#F)', true

'class{ F(){} }().Method?(#F)', true

'Def("C", "class { F(){} }")
C{}.Method?(#F)', true

'Def("C", "class { F(){} }")
C{}().Method?(#F)', true

'x = class{}()
x.F = function (){}
x.Method?(#F)', false

// MethodClass

'c = class { F(){} }
c.MethodClass(#F) is c', true

'c = class { F(){} }
c().MethodClass(#F) is c', true

'Def("C", "class { F(){} }")
C.MethodClass(#F) is C', true

'Def("C", "class { F(){} }")
C().MethodClass(#F) is C', true

'Def("C", "class { F(){} }")
c = C { }
c.MethodClass(#F) is C', true

'Def("C", "class { F(){} }")
c = C { }
c().MethodClass(#F) is C', true

'Def("C", "class { F(){} }")
c = C { F(){} }
c.MethodClass(#F) is c', true

'Def("C", "class { F(){} }")
c = C { F(){} }
c().MethodClass(#F) is c', true

'x = class{}()
x.F = function (){}
x.MethodClass(#F)', false

// GetDefault

'class{}.GetDefault(#X, 123)', 123

'class{X: 123}.GetDefault(#X, 456)', 123

'class{Getter_X(){ 123 }}.GetDefault(#X, 456)', 123

'class{Getter_X(){ 123 }}().GetDefault(#X, 456)', 123

'class{Getter_(member){ 123 }}.GetDefault(#X, 456)', 123

'class{Getter_(member){ 123 }}().GetDefault(#X, 456)', 123

'class{}().GetDefault(#X, 123)', 123

'class{X: 123}().GetDefault(#X, 456)', 123

'x = class{}()
x.A = 123
x.GetDefault(#A, 456)', 123

// Copy
'Def("C", "class { }")
x = C{}()
x.A = 123
x.B = 456
y = x.Copy()
[y.Base?(C), y.A, y.B]', '#(true, 123, 456)'

// Default method

'c = class { Default(@args) { args } }
c.Foo()', '#(Foo)'

'c = class { Default(@args) { args } }
c.Foo(12,34)', '#(Foo,12,34)'

'c = class { Default(@args) { args } }
c.Foo(@#(1,2,a:3))', '#(Foo,1,2,a:3)'

'c = class { Default(@args) { args } }
c.Foo(@+1#(1,2,a:3))', '#(Foo,2,a:3)'

'c = class { Default(@args) { args } }
c().Foo(12,34)', '#(Foo,12,34)'

'c = class { Foo(x,y) { x + y } Default(@args) { args } }
c.Foo(12,34)', '46'

'c = class { Foo(x,y) { x + y } Default(@args) { args } }
c().Foo(12,34)', '46'

'// Default works with call
i = new class { Default(@args) { args } }
i()', '#(Call)'

// Eval

'c = class { X: 123 }
c.Eval({ .X })', 123

'c = class { }
x = c.Eval({ ; })' throws 'no return value'

'c = class { X: 123 }
c.Eval2({ .X })', '#(123)'

'c = class { X: 123 }
c().Eval({ .X })', 123

'c = class { }
x = c().Eval({ ; })' throws 'no return value'

'x = class{}()
x.X = 123
x.Eval({ .X })', 123
`)
