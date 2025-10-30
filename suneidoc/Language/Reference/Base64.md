### Base64

``` suneido
.Encode(string) => base64_string
```

``` suneido
.Decode(base64_string) => string
```

Returns the encoded or decoded string.

For example:

``` suneido
Base64.Encode("hello") => "aGVsbG8="

Base64.Decode("aGVsbG8=") => "hello"
```

**Note**: Also available as [string.Base64Encode](<String/string.Base64Encode.md>) and [string.Base64Decode](<String/string.Base64Decode.md>)

See also:
[RFC2045 MIME Part One: Format of Internet Message Bodies](<http://www.faqs.org/rfcs/rfc2045.html>)