<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Detab

``` suneido
() => string
```

Returns a copy of the string with tab characters replaced by spaces, assuming tab stops every four characters.

For example:

``` suneido
"a\tb".Detab() => "a   b"
```

See also:
[string.Entab](<string.Entab.md>)