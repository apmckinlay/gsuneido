#### string.FindRxLast

``` suneido
(regexp) => number or false
```

Returns the last position in the string where the regexp is matched, or false if it's not found.

For example:

``` suneido
"this is a test".FindRxLast("[aeiou]s") => 11
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
