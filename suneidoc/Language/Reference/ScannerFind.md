### ScannerFind

``` suneido
(text, string) => number
```

Returns the position of string in text, or text.Size() if the string is not found.

Differs from [string.Find](<String/string.Find.md>) because it uses Scanner to search the text a token at a time.

For example:

``` suneido
s = "/* x */ x"
s.Find("x")
    => 3
ScannerFind(s, "x")
    => 8
```

However, the comparison is done to the remainder of the string, rather than to one token, so you can search for a sequence of several tokens.

``` suneido
ScannerFind("x x y", "x y")
    => 2
```

**Note:** This is an exact match (with [string.Prefix?](<String/string.Prefix?.md>)) so whitespace is significant.

See also: [string.Find](<String/string.Find.md>), 
[Scanner](<Scanner.md>)