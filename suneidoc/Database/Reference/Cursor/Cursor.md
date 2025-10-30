<div style="float:right"><span class="builtin">Builtin</span></div>

#### Cursor

``` suneido
Cursor(query [, field: value ...]) => cursor or false
Cursor(query [, field: value ...], block) => value
```

Returns a cursor for a query.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

If a block (or other callable value) is supplied, it is called with the cursor as its argument, the cursor is automatically Close'd when it returns, and the return value from the block is returned. If the block throws an exception, Close is still done. This is a good way to ensure that cursors are closed.

Throws "Cursor: invalid query" if the query is impossible, or not cursor-able.

See also: [transaction.Query](<../Transaction/transaction.Query.md>)