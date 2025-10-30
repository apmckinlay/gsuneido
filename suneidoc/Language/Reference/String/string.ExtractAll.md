#### string.ExtractAll

``` suneido
(pattern) => object
```

Returns all of the matches for a given pattern,
or false if the string doesn't match the pattern.

For example:

``` suneido
"hello world".ExtractAll("(h.*)\s(w.*)") => #("hello world", "hello", "world")
```

**Note:** '+' and '*' in regular expressions
are currently implemented recursively
and are limited to recursing at most 500 times.
This limits the quantity of text than can be Extract'ed.

See also:
[Regular Expressions](<../../Regular Expressions.md>),
[string.Match](<string.Match.md>),
[string.Replace](<string.Replace.md>), 
[string.Extract](<string.Extract.md>)