### QueryAddWhere

``` suneido
(query, where) => query
```

Returns the query with the *where* added, prior to any sort.

For example:

``` suneido
query = "tables sort tablename"
QueryAddWhere(query, "where table = 1")
    => "tables where table = 1 sort tablename"
```