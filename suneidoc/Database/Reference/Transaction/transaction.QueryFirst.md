<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.QueryFirst

``` suneido
(query [, field: value ...]) => record
```

Fetches the first record from the given query. Returns false if the query does not contain any records.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

**Note**: The query must have a sort, otherwise first/last are undefined.


See also:
[Query1](<../Query1.md>),
[QueryEmpty?](<../QueryEmpty?.md>),
[QueryFirst](<../QueryFirst.md>),
[QueryLast](<../QueryLast.md>),
[Query.Strategy1](<../Query/Query.Strategy1.md>),
[transaction.Query1](<transaction.Query1.md>),
[transaction.QueryEmpty?](<transaction.QueryEmpty?.md>),
[transaction.QueryLast](<transaction.QueryLast.md>)
