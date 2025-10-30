### QueryColumns

``` suneido
(query, useCache = false) => list
```

Returns a list of columns for the query. If a table does not exist it returns an empty object.

If useCache is true then results are cached for 60 seconds.

For example:

``` suneido
QueryColumns("tables join columns")
    => #("table", "tablename", "nextfield", "nrows", "totalsize", "column", "field")
```

See also:
[query.Columns](<Query/query.Columns.md>)