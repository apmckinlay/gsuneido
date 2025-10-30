<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.Query

``` suneido
(query [, field: value ...]) => query or true or false
(query [, field: value ...], block) => value
```

Returns a [Query](<../Query.md>) object for the given query.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

If a block (or other callable value) is supplied, it is called with the query as its argument, the query is automatically Close'd when it returns, and the return value from the block is returned. If the block throws an exception, Close is still done. This is a good way to ensure that queries are closed.

``` suneido
Transaction(read:)
    { |t|
    t.Query('tables')
        { |q|
        first = q.Next()
        second = q.Next()
        }
    }
```

Throws "Transaction.Query: invalid query" if the query is impossible.

See also:
[WithQuery](<../WithQuery.md>),
[QueryDo](<../QueryDo.md>),
[Cursor](<../Cursor.md>)