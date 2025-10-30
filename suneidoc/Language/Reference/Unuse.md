<div style="float:right"><span class="builtin">Builtin</span></div>

### Unuse

``` suneido
( string ) => true or false
```

Stop using the named library.  Returns false if the library is not in use.

**Note**: As a side-effect, this will unload all library records from memory. This can have a significant effect on performance.

See also:
[Libraries](<Libraries.md>),
[LibraryTables](<LibraryTables.md>),
[Use](<Use.md>)