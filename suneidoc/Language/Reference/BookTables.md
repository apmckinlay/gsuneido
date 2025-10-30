### BookTables

``` suneido
() => object
```

Returns a list of names of tables which have the fields required by books.
In most cases this will be a list of the available books.

For example:

``` suneido
BookTables()
    => #('mybook', 'suneidoc')
```

See also: [LibraryTables](<LibraryTables.md>)