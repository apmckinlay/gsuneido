<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.FromUtf8

``` suneido
() => string
```

Returns the string converted from UTF-8 to the current code page.

**Warning**: Converting with a Windows code page (anything other than CP.UTF8) may not be portable.


See also:
[string.ToUtf8](<string.ToUtf8.md>),
[MultiByteToWideChar](<../MultiByteToWideChar.md>),
[WideCharToMultiByte](<../WideCharToMultiByte.md>)
