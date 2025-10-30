#### object.Empty?

``` suneido
() => true or false
```

Returns true if the object has no members, i.e. Size() is 0, false if not.

For example:

``` suneido
#().Empty?() => true
#(123).Empty?() => false
#(name: "Fred").Empty?() => false
```

See also: [object.NotEmpty?](<object.NotEmpty?.md>)