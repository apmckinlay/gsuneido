### Query.StripSort

``` suneido
(query) => query
```

Returns the query with the *sort* removed (if there is one).

For example:

``` suneido
query = "tables join columns sort column"
Query.StripSort(query)
    => "tables join columns"
```