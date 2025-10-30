### QueryTotal

``` suneido
(query, field) => number
```

Returns the total for the field in the query.

For example:

``` suneido
QueryTotal("tables", "nrows")
    => 1234
```

See also:
[transaction.QueryTotal](<Transaction/transaction.QueryTotal.md>),
[QueryCount](<QueryCount.md>),
[QueryMax](<QueryMax.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMin](<QueryMin.md>),
[QueryMode](<QueryMode.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),