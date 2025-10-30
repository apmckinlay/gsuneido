<div style="float:right"><span class="builtin">Builtin</span></div>

### MultiByteToWideChar

``` suneido
(string, cp = CP.ACP) => string
```

The returned string uses two bytes to represent each character. The result is terminated with two 0 bytes (i.e. a wide nul)

Only handles 1252 and CP.UTF8. ACP and THREAD_ACP are treated as 1252.

**Warning**: Converting with a Windows code page (anything other than CP.UTF8) may not be portable.


See also:
[string.FromUtf8](<String/string.FromUtf8.md>),
[string.ToUtf8](<String/string.ToUtf8.md>),
[WideCharToMultiByte](<WideCharToMultiByte.md>)
