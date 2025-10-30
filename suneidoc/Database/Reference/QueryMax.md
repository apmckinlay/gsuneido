### QueryMax

``` suneido
(query, field, default = false) => value
```

Returns the maximum value of the field in the query, or the supplied default if there are no records in the query.

For example:

``` suneido
QueryMax("tables", "table")
    => 123
```

See also:
[transaction.QueryMax](<Transaction/transaction.QueryMax.md>),
[QueryCount](<QueryCount.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMin](<QueryMin.md>),
[QueryMode](<QueryMode.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),
[QueryTotal](<QueryTotal.md>)