### transaction.QueryCount

``` suneido
(query) => number
```

Returns the number of records in a query.

For example:

``` suneido
t.QueryCount("tables where name =~ 'lib'")
    => 5
```

See also:
[QueryCount](<../QueryCount.md>),
[transaction.QueryList](<transaction.QueryList.md>),
[transaction.QueryMean](<transaction.QueryMean.md>),
[transaction.QueryMin](<transaction.QueryMin.md>),
[transaction.QueryMax](<transaction.QueryMax.md>),
[transaction.QueryTotal](<transaction.QueryTotal.md>)