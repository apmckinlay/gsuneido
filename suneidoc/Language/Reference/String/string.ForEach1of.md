#### string.ForEach1of

``` suneido
(chars, block) => number
```

Calls block with each position that is one of chars.

For example:

``` suneido
"hello world".ForEach1of("abcde")
	{ Print(it) }
=>	1
	10
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
