### QueryStrategy

``` suneido
(query, formatted = false) => string
```

Returns the *strategy* for the query (using [query.Strategy](<Query/query.Strategy.md>)).

For example:

``` suneido
QueryStrategy("tables join columns")
    => "(tables^(table)) JOIN 1:n on (table) (columns^(table,column))"
```

See also: [Query.Strategy1](<Query/Query.Strategy1.md>)