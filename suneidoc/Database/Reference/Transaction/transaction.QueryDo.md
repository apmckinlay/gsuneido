<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.QueryDo

``` suneido
(request [, field: value ...]) => count
```

Executes an insert, update, and delete request and returns the number of records processed.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

(This is only applicable for delete requests.)

See also: [QueryDo](<../QueryDo.md>)