### UrlDecode

``` suneido
(string) => string
```

Converts '+' to space, 
and %hh (where h is a hex digit) to the corresponding ASCII character.

Useful for
[HttpServer](<HttpServer.md>)
functions.

See also:
[UrlDecodeValues](<UrlDecodeValues.md>),
[UrlEncode](<UrlEncode.md>)