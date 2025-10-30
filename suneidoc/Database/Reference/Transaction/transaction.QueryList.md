### QueryList

``` suneido
(query, field) => number
```

Returns a list of values found in the specified field from the query.

QueryList removes any sort, appends "summarize list", and then executes the query.

**Warning:** The list is accumulated in memory so there are limits to it's size.

For example:

``` suneido
t.QueryList("columns where table = 1", "column")
    => #("nextfield", "nrows", "table", "tablename", "totalsize")
```

See also:
[QueryList](<../QueryList.md>),
[transaction.QueryCount](<transaction.QueryCount.md>),
[transaction.QueryMax](<transaction.QueryMax.md>),
[transaction.QueryMean](<transaction.QueryMean.md>),
[transaction.QueryMin](<transaction.QueryMin.md>),
[transaction.QueryTotal](<transaction.QueryTotal.md>)