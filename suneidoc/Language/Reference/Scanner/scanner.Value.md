<div style="float:right"><span class="builtin">Builtin</span></div>

#### scanner.Value

``` suneido
() => string
```

Returns the value of the current token.

[scanner.Text](<scanner.Text.md>) and scanner.Value are almost the same, except for strings, where Value will not include the quotes and will handle escapes (e.g. converting \t to a tab character) whereas Text will be the raw source.