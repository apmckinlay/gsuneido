<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.ToUtf8

``` suneido
() => string
```

Returns the string converted to UTF-8 based on the current code page.

**Warning**: Converting with a Windows code page (anything other than CP.UTF8) may not be portable.


See also:
[string.FromUtf8](<string.FromUtf8.md>),
[MultiByteToWideChar](<../MultiByteToWideChar.md>),
[WideCharToMultiByte](<../WideCharToMultiByte.md>)
