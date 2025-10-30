### transaction.QueryMax

``` suneido
(query, field, default = false) => value
```

Returns the maximum value of the field in the query, or the supplied default if there are no records in the query.

For example:

``` suneido
t.QueryMax("tables", "table")
    => 123
```

See also:
[QueryMax](<../QueryMax.md>),
[transaction.QueryMean](<transaction.QueryMean.md>),
[transaction.QueryMin](<transaction.QueryMin.md>),
[transaction.QueryCount](<transaction.QueryCount.md>),
[transaction.QueryTotal](<transaction.QueryTotal.md>)