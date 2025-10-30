<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Has?

``` suneido
(value) => true or false
```

Returns true if the object contains the specified value, false if not.

For example:

``` suneido
#(123, "abc", name: "Fred", age: 23).Has?("Fred") => true
#(123, "abc", name: "Fred", age: 23).Has?(123) => true
#(123, "abc", name: "Fred", age: 23).Has?(3) => false
#(123, "abc", name: "Fred", age: 23).Has?("name") => false
```


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.HasIf?](<object.HasIf?.md>)
