<div style="float:right"><span class="builtin">Builtin</span></div>

### Database.Check

``` suneido
() => "" or error string
```

Check the integrity of the database while the server is running.

Equivalent to the `-check` [command line option](<../../../Introduction/Command Line Options.md>)

Returns "" if ok, otherwise a multi-line error string (the same as the output from -check).