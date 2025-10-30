### KeyQuery

``` suneido
(query, record)
```

Build a query for a record by adding a "where" to the query using the [ShortestKey](<ShortestKey.md>). This query can be used to update or delete the record.

**Note:** The record must contain the fields for the key (with the same names).

If the query has an empty key (no fields), then the query is returned unaltered, since there can only be at most a single record in this case.

For example:

``` suneido
KeyQuery('columns', #(table: 'stdlib', column: 'num'))

    => 'columns where table = "stdlib" and column = "num"'
```