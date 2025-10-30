### UrlEncode

``` suneido
(string) => string
```

Converts spaces to '+', and other special characters to %hh (where h is a hex digit).

**Note:** to allow actual URLs to be encoded, this does not encode URL syntax characters (";", "/", "?", ":", "@", "=", "#" and "&")

See also:
[UrlDecode](<UrlDecode.md>)