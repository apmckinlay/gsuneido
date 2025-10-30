<div style="float:right"><span class="builtin">Builtin</span></div>

### Query1

``` suneido
(query [, field: value ...]) => record
```

Returns the single record from the given query. Returns false if the query does not contain any records. Throws an exception if there is more than one record. Executes the query in a standalone read-only transaction.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

Calls with just a table name in the first argument, and named arguments for the rest of the query are often faster because they bypass query parsing and use simpler streamlined optimization that does not look at the data.


See also:
[QueryEmpty?](<QueryEmpty?.md>),
[QueryFirst](<QueryFirst.md>),
[QueryLast](<QueryLast.md>),
[Query.Strategy1](<Query/Query.Strategy1.md>),
[transaction.Query1](<Transaction/transaction.Query1.md>),
[transaction.QueryEmpty?](<Transaction/transaction.QueryEmpty?.md>),
[transaction.QueryFirst](<Transaction/transaction.QueryFirst.md>),
[transaction.QueryLast](<Transaction/transaction.QueryLast.md>)
