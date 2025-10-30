<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Upper?

``` suneido
( ) => true or false
```

Returns true if all letters in the string are upper case (capitals), else returns false.

For example:

``` suneido
"".Upper?()
    => false
"123".Upper?()
    => false
"Hello".Upper?()
    => false
"HELLO WORLD 123".Upper?()
    => true
```

Use [string.Lower?](<string.Lower?.md>) to check for lower case.


See also:
[string.Alpha?](<string.Alpha?.md>),
[string.AlphaNum?](<string.AlphaNum?.md>),
[string.Blank?](<string.Blank?.md>),
[string.Capitalized?](<string.Capitalized?.md>),
[string.Has1of?](<string.Has1of?.md>),
[string.Has?](<string.Has?.md>),
[string.Lower?](<string.Lower?.md>),
[string.Number?](<string.Number?.md>),
[string.Numeric?](<string.Numeric?.md>),
[string.Prefix?](<string.Prefix?.md>),
[string.Suffix?](<string.Suffix?.md>),
[string.White?](<string.White?.md>)
