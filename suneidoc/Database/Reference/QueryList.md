### QueryList

``` suneido
(query, field) => object
```

Returns a list of values found in the specified field from the query.

QueryList removes any sort, appends "summarize list", and then executes the query.

**Warning:** The list is accumulated in memory so there are limits to its size.

For example:

``` suneido
QueryList("columns where table = 1", "column")
    => #("nextfield", "nrows", "table", "tablename", "totalsize")
```

See also:
[transaction.QueryList](<Transaction/transaction.QueryList.md>),
[QueryCount](<QueryCount.md>),
[QueryMax](<QueryMax.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMin](<QueryMin.md>),
[QueryMode](<QueryMode.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),
[QueryTotal](<QueryTotal.md>)