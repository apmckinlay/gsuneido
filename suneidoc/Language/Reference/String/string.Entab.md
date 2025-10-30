<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Entab

``` suneido
() => string
```

Returns a copy of the string with <u>leading</u> spaces replaced by tabs, assuming tab stops every four characters.

Also removes trailing spaces or tabs from lines

For example:

``` suneido
"    hello ".Entab() => "\thello"
```

See also:
[string.Detab](<string.Detab.md>)