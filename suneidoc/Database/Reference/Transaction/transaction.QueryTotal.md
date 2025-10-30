### transaction.QueryTotal

``` suneido
(query, field) => number
```

Returns the total for the field in the query.

For example:

``` suneido
t.QueryTotal("tables", "nrows")
    => 1234
```

See also:
[QueryTotal](<../QueryTotal.md>),
[transaction.QueryMean](<transaction.QueryMean.md>),
[transaction.QueryMin](<transaction.QueryMin.md>),
[transaction.QueryMax](<transaction.QueryMax.md>),
[transaction.QueryCount](<transaction.QueryCount.md>),