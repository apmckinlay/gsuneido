### transaction.QueryMin

``` suneido
(query, field, default = false) => value
```

Returns the minimum value of the field in the query, or the supplied default if there are no records in the query.

For example:

``` suneido
t.QueryMin("tables", "table")
    => 0
```

See also:
[QueryMin](<../QueryMin.md>),
[transaction.QueryMax](<transaction.QueryMax.md>),
[transaction.QueryCount](<transaction.QueryCount.md>),
[transaction.QueryTotal](<transaction.QueryTotal.md>)