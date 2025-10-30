### QueryCount

``` suneido
(query) => number
```

Returns the number of records in a query.

QueryCount removes any sort, appends "summarize count", and then executes the query.

For example:

``` suneido
QueryCount("tables where name =~ 'lib'")
    => 5
```

See also:
[transaction.QueryCount](<Transaction/transaction.QueryCount.md>),
[QueryMax](<QueryMax.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMin](<QueryMin.md>),
[QueryMode](<QueryMode.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),
[QueryTotal](<QueryTotal.md>)