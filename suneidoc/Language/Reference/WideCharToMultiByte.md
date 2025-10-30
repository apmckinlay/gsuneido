<div style="float:right"><span class="builtin">Builtin</span></div>

### WideCharToMultiByte

``` suneido
(string, cp = CP.ACP) => string
```

The input string must be terminated with two 0 bytes (i.e. a wide nul)

The returned string does not have an explicit nul terminator. But as with any Suneido string, it will be null terminated if passed to a dll call.

Only handles 1252 and CP.UTF8. ACP and THREAD_ACP are treated as 1252.

**Warning**: Converting with a Windows code page (anything other than CP.UTF8) may not be portable.


See also:
[string.FromUtf8](<String/string.FromUtf8.md>),
[string.ToUtf8](<String/string.ToUtf8.md>),
[MultiByteToWideChar](<MultiByteToWideChar.md>)
