<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Find

``` suneido
(string, pos = 0) => number
```

Returns the first position greater than or equal to **pos** where **string** is found, or Size() if not found.

For example:

``` suneido
"this is a test".Find("is") => 2
"hello world".Find("o") => 4
"hello world".Find("o", 5) => 7
"hello world".Find("x") => 11
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
