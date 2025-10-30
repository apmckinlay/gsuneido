<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Match

``` suneido
(pattern, pos = false, prev = false) => Object(Object(i, n), ...
```

Returns information about the parts of the string matched by the pattern, or false if there is no match.

If **prev** is true it searches in reverse.

**pos** can be specified to start the search somewhere other than the ends of the string.

The result object will contain one sub-object for the entire match, plus additional sub-objects for any parenthesized sub-patterns.

For example:

``` suneido
"hello world".Match("w(..)ld") => #((6,5),(7,2))
"hello world".Match("o", prev:) => #((7,1))
```

See also:
[Regular Expressions](<../../Regular Expressions.md>),
[string.Extract](<string.Extract.md>),
[string.Replace](<string.Replace.md>)


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>)
