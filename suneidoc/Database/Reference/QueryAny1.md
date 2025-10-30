### QueryAny1

``` suneido
(query, [, field ...]) => object
```

Fetches a random record from the given query. Returns false if the query does not contain any records, or an object with the fields passed in as members and their values. Executes the query in a standalone readonly transaction.

Useful when a group of records has the same value so you only need one of them and it doesn't matter which. Does not use a sort.

Can use **Suneido.ValidateQueryAny1? = true** to turn on validation that will check a max of 10 records to ensure they have the same values for the list of fields

For example:

``` suneido
group = QueryAny1("stdlib where name =~ 'File'", "group").group
```

See also:
[Query1](<Query1.md>),
[QueryLast](<QueryLast.md>),
[QueryFirst](<QueryFirst.md>), 
[transaction.QueryAny1](<Transaction/transaction.QueryAny1.md>)