### QueryMode

``` suneido
(query, field) => object
```

Returns the mode (a list of values) for the field in the query.

For example:

``` suneido
QueryMode("tables", "nextfield")
    => #(7)
```

See also:
[QueryCount](<QueryCount.md>),
[QueryMax](<QueryMax.md>),
[QueryMean](<QueryMean.md>),
[QueryMeanDev](<QueryMeanDev.md>),
[QueryMedian](<QueryMedian.md>),
[QueryMin](<QueryMin.md>),
[QueryRange](<QueryRange.md>),
[QueryStdDev](<QueryStdDev.md>),
[QueryTotal](<QueryTotal.md>)