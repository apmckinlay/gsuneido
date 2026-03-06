// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

var _ = Register("records 1", `
'Record?(Record())'
'Record?(#{})'
'Record?([])'

'r=[a: 1, b: 2];r.Clear();r', '[]'
'Object().Clear()' throws "not found"
'r=[a: 1, b: 2];r.Delete(all:);r', '[]'
`)

var _ = Register("records observers", `
// multiple observers
'r = Record()
r.Observer({|member| o1 = member })
r.Observer({|member| o2 = member })
r.foo = 123
[o1, o2]', '#(foo,foo)'

// RemoveObserver
'r = Record()
r.Observer(b1 = {|member| o1 = member })
r.Observer({|member| o2 = member })
r.foo = 123
r.RemoveObserver(b1)
r.bar = 456
[o1, o2]', '#(foo,bar)'

// observer called as method
'r = Record()
r.Observer({|member| o = member $ "=" $ this[member] })
r.foo = 123
o', '"foo=123"'

// recursive with different members allowed
'r = Record()
r.Observer({|member| o = member; .bar = 456 })
r.foo = 123
o', 'bar'

// PreSet does not call observers
'r = Record()
o = false
r.Observer({|member| o = member })
r.PreSet(#foo, 123)
[o,r.foo]', '#(false,123)'

// queued observers
'r = Record()
r.Observer({|member| o1 = member; r.bar = 1 })
r.Observer({|member| o2 = member })
r.foo = 123
[o1,o2]', '#(bar, foo)'

// observer not called for assignment of equal value
'r = Record(x: 123)
o = ""
r.Observer({|member| o = member })
r.x = 123
o', '""'

'v = #(123)
r = Record(x: v)
o = ""
r.Observer({|member| o = member })
r.x = v
o', '""'

// gSuneido only behavior
// 'r = Record(x: Object(123))
// o = ""
// r.Observer({|member| o = member })
// r.x = Object(123)
// o', '"x"'
`)

var _ = Register("records rules", `
// _lower! with no value
'r = Record()
r.foo_lower!', '""'

// _lower! with non-string
'Record(foo: true).foo_lower!', 'true'

// _lower!
'Record(foo: "Hello World").foo_lower!', '"hello world"'

// simple implicit rule, no dependencies
'Def("Rule_foo", function() { 123 })
Record().foo', 123

// removing default disables implicit rules
'Def("Rule_foo", function() { 123 })
r = Record()
r.Set_default()
r.foo' throws "member not found"

// rule saves result
'Def("Rule_foo", function() { 123 })
r = Record()
r.foo
r', '[foo: 123]'

// rule does not save result if read-only
'Def("Rule_foo", function() { 123 })
r = #{}
r.foo
r', '#{}'

// attached rule
'r = Record()
r.AttachRule("foo", function(){ 123 })
r.foo', 123

// attached rule still applies even if default removed to disable implicit rules
'r = Record()
r.AttachRule("foo", function(){ 123 })
r.Set_default()
r.foo', 123

// attached rule overrides implicit rule
'Def("Rule_foo", function() { 123 })
r = Record()
r.AttachRule("foo", function(){ 456 })
r.foo', 456

// rule exception is wrapped
'r = Record()
r.AttachRule("foo", function(){ throw "error" })
r.foo' throws "foo)"

// rule doesn"t trigger itself
'r = Record()
r.AttachRule("foo", function(){ .foo; 123 })
r.foo', 123

// rule dependencies tracked
'Def("Rule_foo", function() { .bar $ .baz })
r = Record()
r.foo
r.GetDeps(#foo).Split(",").Sort!()', '#(bar,baz)'

// rule SetDeps/GetDeps roundtrip
'r = Record()
r.SetDeps(#foo, "one, two,three ")
r.GetDeps(#foo).Split(",").Sort!()', '#(one,three,two)'

// rule result cached, only called when missing/invalid
'n = 0
r = Record()
r.AttachRule("foo", { n++; .bar })
r.foo
r.foo
n' 1

// rule result cached, called again when dependency changes
'n = 0
r = Record()
r.AttachRule("foo", { n++; .bar })
r.foo
r.foo
r.bar = 123
x = r.foo
[n,x]', '#(2,123)'

// rule not invalidated if value does not change
'n = 0
r = Record(bar: 123)
r.AttachRule("foo", { n++; .bar })
r.foo
r.foo
r.bar = 123
x = r.foo
[n,x]', '#(1,123)'

// rule result cached, called again when invalidated
'n = 0
r = Record(bar: 123)
r.AttachRule("foo", { n++; .bar })
r.foo
r.foo
r.Invalidate(#foo)
x = r.foo
[n,x]', '#(2,123)'
`)
