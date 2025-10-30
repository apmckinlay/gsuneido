<div style="float:right"><span class="builtin">Builtin</span></div>

### QueryFirst

``` suneido
(query [, field: value ...]) => record or false
```

Fetches the first record from the given query. Returns false if the query does not contain any records. Executes the query in a standalone read-only transaction.

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
[Query1](<Query1.md>),
[QueryEmpty?](<QueryEmpty?.md>),
[QueryLast](<QueryLast.md>),
[Query.Strategy1](<Query/Query.Strategy1.md>),
[transaction.Query1](<Transaction/transaction.Query1.md>),
[transaction.QueryEmpty?](<Transaction/transaction.QueryEmpty?.md>),
[transaction.QueryFirst](<Transaction/transaction.QueryFirst.md>),
[transaction.QueryLast](<Transaction/transaction.QueryLast.md>)
