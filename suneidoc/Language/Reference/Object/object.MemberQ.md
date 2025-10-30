<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Member?

``` suneido
(member) => true or false
```

Returns true if the member exists, false otherwise.

For example:

``` suneido
#(12, 34, a: 56, b: 78).Member?(1) => true
#(12, 34, a: 56, b: 78).Member?(2) => false
#(12, 34, a: 56, b: 78).Member?(56) => false
#(12, 34, a: 56, b: 78).Member?("a") => true
#(12, 34, a: 56, b: 78).Member?("c") => false
```