<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.FindLast1of

``` suneido
(chars, pos = Size()) => number or false
```

Returns the last position less than or equal to **pos** where a character from **chars** is found, or false if not found.

**chars** is a list of characters.

If the first character is '^' then the set of characters will be negated. 
i.e. it will match everything except those characters. A literal '^' can be included as long as it is not the first character. e.g. "abc^"

A range of characters can be abbreviated with '-'.
For example 'a-zA-Z'. 
A literal '-' can be included as the first or last character e.g. "-+*/" or "*/+="

Note: negation and ranges were added 2024-10-25

For example:

``` suneido
"this is a test".FindLast1of("si") => 12
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEachMatch](<string.ForEachMatch.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
