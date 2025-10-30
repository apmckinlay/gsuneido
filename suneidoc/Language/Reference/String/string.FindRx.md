#### string.FindRx

``` suneido
(regexp) => number
```

Returns the first position within the string where the regexp is matched, or Size() if it's not found.

For example:

``` suneido
"this is a test".FindRx("t[aeiou]") => 10
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
