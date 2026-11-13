### Global

``` suneido
(name) => value
```

Returns the value for a global (capitalized) name.

Also handles static public members e.g. 'LC.TIME'

Throws an exception if the name is not a valid global name.

Safer than [string.Eval](<String/string.Eval.md>) because it will not accept arbitrary code.