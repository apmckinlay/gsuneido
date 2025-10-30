<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Numeric?

``` suneido
( ) => true or false
```

Return true if all characters in the string are numeric (digits) and there is at least one character, false otherwise.

For example:

``` suneido
"".Numeric?()
    => false
"123".Numeric?()
    => true
"123abc".Numeric?()
    => false
```

To determine if a string is a number use [string.Number?](<string.Number?.md>)


See also:
[string.Alpha?](<string.Alpha?.md>),
[string.AlphaNum?](<string.AlphaNum?.md>),
[string.Blank?](<string.Blank?.md>),
[string.Capitalized?](<string.Capitalized?.md>),
[string.Has1of?](<string.Has1of?.md>),
[string.Has?](<string.Has?.md>),
[string.Lower?](<string.Lower?.md>),
[string.Number?](<string.Number?.md>),
[string.Prefix?](<string.Prefix?.md>),
[string.Suffix?](<string.Suffix?.md>),
[string.Upper?](<string.Upper?.md>),
[string.White?](<string.White?.md>)
