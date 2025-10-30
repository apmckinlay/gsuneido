### Global

``` suneido
(name) => value
```

Returns the value for a global (capitalized) name.

Also handles static public members e.g. 'LC.TIME'

Throws an exception if the name is not a valid global name.

Uses [string.Eval](<String/string.Eval.md>) but is safer because it will not accept arbitrary code.