### QueryMin

``` suneido
(query, field, default = false) => value
```

Returns the minimum value of the field in the query, or the supplied default if there are no records in the query.

For example:

``` suneido
QueryMin("tables", "table")
    => 0
```

See also:
[transaction.QueryMin](<Transaction/transaction.QueryMin.md>),
[QueryCount](<QueryCount.md>),
[QueryMax](<QueryMax.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMode](<QueryMode.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),
[QueryTotal](<QueryTotal.md>)