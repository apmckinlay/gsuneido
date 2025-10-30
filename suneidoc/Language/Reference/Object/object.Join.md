<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Join

``` suneido
(separator = "") => string
```

Concatenates the un-named members, with the specified separator string between them.

For example:

``` suneido
Object("one", "two", "three").Join(", ") => "one, two, three"
```

See also: 
[Join](<../Join.md>),
[string.Split](<../String/string.Split.md>)