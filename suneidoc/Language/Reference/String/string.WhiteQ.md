#### string.White?

``` suneido
( ) => true or false
```

Return true if all characters in the string are whitespace (blank, tab, return, or linefeed) and there is at least one character, false otherwise.

Similar to [string.Blank?](<string.Blank?.md>) except that White? returns true for "" (empty string).

For example:

``` suneido
"".White?()
    => false
"a b c".White?()
    => false
"   ".White?()
    => true
```


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
[string.Upper?](<string.Upper?.md>)
