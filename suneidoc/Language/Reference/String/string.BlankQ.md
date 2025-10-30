#### string.Blank?

``` suneido
( ) => true or false
```

Return true if all characters in the string are whitespace (blank, tab, return, or linefeed), false otherwise.

Similar to [string.White?](<string.White?.md>) except that Blank? returns true for "" (empty string).

For example:

``` suneido
"".Blank?()
    => true
"a b c".Blank?()
    => false
"   ".Blank?()
    => true
```


See also:
[string.Alpha?](<string.Alpha?.md>),
[string.AlphaNum?](<string.AlphaNum?.md>),
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
