<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Lower?

``` suneido
() => true or false
```

Returns true if all letters are lower case (not capitals), else returns false.

For example:

``` suneido
"".Lower?()
    => false
"123".Lower?()
    => false
"Hello".Lower?()
    => false
"hello world 123".Lower?()
    => true
```

Use [string.Upper?](<string.Upper?.md>) to check for lower case.


See also:
[string.Alpha?](<string.Alpha?.md>),
[string.AlphaNum?](<string.AlphaNum?.md>),
[string.Blank?](<string.Blank?.md>),
[string.Capitalized?](<string.Capitalized?.md>),
[string.Has1of?](<string.Has1of?.md>),
[string.Has?](<string.Has?.md>),
[string.Number?](<string.Number?.md>),
[string.Numeric?](<string.Numeric?.md>),
[string.Prefix?](<string.Prefix?.md>),
[string.Suffix?](<string.Suffix?.md>),
[string.Upper?](<string.Upper?.md>),
[string.White?](<string.White?.md>)
