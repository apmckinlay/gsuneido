#### object.MembersIf

``` suneido
(callable) => object
```

Returns a list of member names of an object for which callable returns True.

callable is called with the member name as its single argument.

callable can be anything that can be called, e.g. block, function, class, or instance.

For example:

``` suneido
#(one: 11, two: 22, three: 33).MembersIf({|m| m =~ 't'})  => #("two", "three")
```

**Note:** Named members are not searched in any particular order.

See also:
[object.Members](<object.Members.md>),
[object.FindIf](<object.FindIf.md>)