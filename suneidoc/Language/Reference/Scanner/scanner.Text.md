<div style="float:right"><span class="builtin">Builtin</span></div>

#### scanner.Text

``` suneido
() => string
```

Returns the text of the current token. (The same as what scanner.Next() returned.)

scanner.Text and [scanner.Value](<scanner.Value.md>) are almost the same, except for strings, where Value will not include the quotes and will handle escapes (e.g. converting \t to a tab character) whereas Text will be the raw source.