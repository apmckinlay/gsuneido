### Schema

``` suneido
( tablename ) => string
```

Returns the database table specification in the same format as required by the table creation requests.

For example:

``` suneido
Schema('tables') =>

"tables
    (nextfield, nrows, table, tablename, totalsize)
    key(table)
    key(tablename)"
```