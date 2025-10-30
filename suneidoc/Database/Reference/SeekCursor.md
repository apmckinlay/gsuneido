### SeekCursor

``` suneido
(query)
```
`Next(tran) => record or false`
: Returns the next row in the query, or false if there are no more rows.

`Prev(tran) => record or false`
: Returns the previous row in the query, or false if there are no more rows.

`Rewind()`
: Rewinds the query so that Next will get the first row or Prev will get the last row.

`Seek(field, prefix)`
: After calling Seek, Next will return the first record greater then or equal to the prefix,
and Prev will return the first record less then the prefix.

`Fields() => object`
: Returns a list of fields.

`Columns() => object`
: Returns a list of the columns.

`Keys() => object`
: Returns a list of the keys.

`Output(transaction, object) => true or false`
: Returns true if the output succeeds, false otherwise.

See also: [Cursor](<Cursor.md>),
[SeekQuery](<SeekQuery.md>)