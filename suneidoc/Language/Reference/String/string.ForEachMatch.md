#### string.ForEachMatch

``` suneido
(regex, block)
```

Calls the block for each non-overlapping match in the string. Passes the result of [string.Match](<string.Match.md>) to the block

For example:

``` suneido
"ab cab abc def".ForEachMatch("ab") {|m| Print(m) }

=>  #(#(0, 2))
    #(#(4, 2))
    #(#(7, 2))
```


See also:
[string.FindLast](<string.FindLast.md>),
[string.Find1of](<string.Find1of.md>),
[string.FindLast1of](<string.FindLast1of.md>),
[string.FindRx](<string.FindRx.md>),
[string.FindRxLast](<string.FindRxLast.md>),
[string.ForEach1of](<string.ForEach1of.md>),
[string.Has?](<string.Has?.md>),
[string.Match](<string.Match.md>)
