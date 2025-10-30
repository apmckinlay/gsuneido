<div style="float:right"><span class="builtin">Builtin</span></div>

### Function?

``` suneido
(value) => true or false
```

Returns true if the value is a function, method, or block, false otherwise.

For example:

``` suneido
if not Function?(x)
    throw "x must be a function"
```


See also:
[Type](<Type.md>),
[Boolean?](<Boolean?.md>),
[Class?](<Class?.md>),
[Date?](<Date?.md>),
[Instance?](<Instance?.md>),
[Number?](<Number?.md>),
[Object?](<Object?.md>),
[Record?](<../../Database/Reference/Record?.md>),
[String?](<String?.md>)
