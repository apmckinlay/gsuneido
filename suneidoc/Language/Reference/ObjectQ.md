<div style="float:right"><span class="builtin">Builtin</span></div>

### Object?

``` suneido
(value) => true or false
```

Returns true if the value is an object or a record, false otherwise.

**Note**: Previously (before 2019-03-20), Object? also returned true for instances of classes.

For example:

``` suneido
if not Object?(x)
    throw "x must be a object"
```


See also:
[Type](<Type.md>),
[Boolean?](<Boolean?.md>),
[Class?](<Class?.md>),
[Date?](<Date?.md>),
[Function?](<Function?.md>),
[Instance?](<Instance?.md>),
[Number?](<Number?.md>),
[Record?](<../../Database/Reference/Record?.md>),
[String?](<String?.md>)
