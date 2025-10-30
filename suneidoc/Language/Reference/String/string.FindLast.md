<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.FindLast

``` suneido
(string, pos = Size()) => number or false
```

Returns the last position less than or equal to **pos** where **string** is found, or false if not found.

For example:

``` suneido
"this is a test".FindLast("is") => 5
"this is a test".FindLast("is", 4) => 2
```


See also:
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
