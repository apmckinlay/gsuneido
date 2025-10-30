<div style="float:right"><span class="builtin">Builtin</span></div>

### Instance?

``` suneido
(value) => true or false
```

Returns true if the value is an instance of a class, false otherwise.

For example:

``` suneido
c = class{}
instance = c()
Instance?(instance)
    => true
```


See also:
[Type](<Type.md>),
[Boolean?](<Boolean?.md>),
[Class?](<Class?.md>),
[Date?](<Date?.md>),
[Function?](<Function?.md>),
[Number?](<Number?.md>),
[Object?](<Object?.md>),
[Record?](<../../Database/Reference/Record?.md>),
[String?](<String?.md>)
