### Query.GetSort

``` suneido
(query) => string
```

Returns the sort fields from the query.

For example:

``` suneido
query = "tables join columns sort table, column"
Query.GetSort(query)
    => "table,column"
```