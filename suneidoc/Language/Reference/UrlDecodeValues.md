### UrlDecodeValues

``` suneido
(string) => object
```

Calls UrlDecode on the string,
and then Split's it on '&'.
If pieces contain an '=' they are stored as named members,
otherwise they are stored as un-named members.

For example:

``` suneido
UrlDecodeValues("hello&world&age=25&text=hello+there")
	=> #("hello", "world", age: 25, text: "hello there")
```

Useful for interpreting query strings in
[HttpServer](<HttpServer.md>)
functions.

See also:
[UrlDecode](<UrlDecode.md>)