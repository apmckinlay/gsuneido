<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Prefix?

``` suneido
(string, pos = 0) => true or false
```

Returns true if the substring starting at the specified position starts with the supplied string, false otherwise.

For example:

``` suneido
"hello world".Prefix?("he") => true
"hello world".Prefix?("world") => false
"hello world".Prefix?("world", 6) => true
```

Use [string.Suffix?](<string.Suffix?.md>) for <u>ends</u> with


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
[string.Suffix?](<string.Suffix?.md>),
[string.Upper?](<string.Upper?.md>),
[string.White?](<string.White?.md>), <a href="/suneidoc/Language/Reference/String/string.RemovePrefix">string.RemovePrefix</a>
