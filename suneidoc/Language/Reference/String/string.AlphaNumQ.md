<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.AlphaNum?

``` suneido
( ) => true or false
```

Return true if all characters in the string are alphanumeric (letters or digits) and there is at least one character, false otherwise.

For example:

``` suneido
"".AlphaNum?()
    => false
"abc123".AlphaNum?()
    => true
"abc123%^*".AlphaNum?()
    => false
```


See also:
[string.Alpha?](<string.Alpha?.md>),
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
