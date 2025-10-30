<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Size

``` suneido
() or (list:) or (named:) => number
```

Returns the number of members in the object, or if list: is specified, the number of un-named members in the object i.e. the initial list.

For example:

``` suneido
#().Size() => 0
#(123, "abc", name: "Fred", age: 25).Size() => 4
#(123, "abc", name: "Fred", age: 25).Size(list:) => 2
#(123, "abc", name: "Fred", age: 25).Size(named:) => 2
```

**Note:** The list: members are the consecutive integer members starting at zero. If a number is "missing" then subsequent numbers will go into the named: members. For example:

``` suneido
ob = Object(11, 22, 33, 44)
ob.Size(list:)
    => 4
ob.Erase(2)
    => #(11, 22, 3: 44)
ob.Size(list:)
    => 2
```