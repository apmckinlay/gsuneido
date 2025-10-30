<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Alpha?

``` suneido
( ) => true or false
```

Return true if all characters in the string are alphabetic (letters) and there is at least one character, false otherwise.

For example:

``` suneido
"".Alpha?()
    => false
"abc".Alpha?()
    => true
"abc123".Alpha?()
    => false
```


See also:
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
[string.Upper?](<string.Upper?.md>),
[string.White?](<string.White?.md>)
