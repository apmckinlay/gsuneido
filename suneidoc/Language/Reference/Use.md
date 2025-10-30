<div style="float:right"><span class="builtin">Builtin</span></div>

### Use

``` suneido
( string ) => true or false
```

Use the named library.  Returns false if the library is already in use.

**Note**: As a side-effect, this will unload all library records from memory. This can have a significant effect on performance.

See also:
[Libraries](<Libraries.md>),
[LibraryTables](<LibraryTables.md>),
[Unuse](<Unuse.md>)