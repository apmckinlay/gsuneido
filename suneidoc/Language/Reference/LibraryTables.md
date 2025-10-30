### LibraryTables

``` suneido
() => object
```

Returns a list of names of tables which have the fields required by libraries.
In most cases this will be a list of the available libraries.

For example:

``` suneido
LibraryTables()
    => #('stdlib', 'apmlib', 'tmplib')
```

See also:
[Libraries](<Libraries.md>),
[Use](<Use.md>),
[BookTables](<BookTables.md>)