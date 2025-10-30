<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Suffix?

``` suneido
(string) => true or false
```

Returns true if the string ends with the supplied string, false otherwise.

For example:

``` suneido
"hello world".Suffix?("he") => false
"hello world".Suffix?("world") => true
```

Use [string.Prefix?](<string.Prefix?.md>) for <u>starts</u> with


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
[string.Upper?](<string.Upper?.md>),
[string.White?](<string.White?.md>), <a href="/suneidoc/Language/Reference/String/string.RemoveSuffix">string.RemoveSuffix</a>
