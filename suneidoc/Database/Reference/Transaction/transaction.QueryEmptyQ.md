#### transaction.QueryEmpty?

``` suneido
(query [, field: value ...]) => true or false
```

Returns true if the query contains no records, otherwise returns false.

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
[Query1](<../Query1.md>),
[QueryEmpty?](<../QueryEmpty?.md>),
[QueryFirst](<../QueryFirst.md>),
[QueryLast](<../QueryLast.md>),
[Query.Strategy1](<../Query/Query.Strategy1.md>),
[transaction.Query1](<transaction.Query1.md>),
[transaction.QueryFirst](<transaction.QueryFirst.md>),
[transaction.QueryLast](<transaction.QueryLast.md>)
