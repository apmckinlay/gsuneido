### SeekQuery

``` suneido
(query)
```
`Next() => record or false`
: Returns the next row in the query, or false if there are no more rows.

`Prev() => record or false`
: Returns the previous row in the query, or false if there are no more rows.

`Rewind()`
: Rewinds the query so that Next will get the first row or Prev will get the last row.

`Seek(field, prefix) => record count or 0`
: Returns the number of records that come before the first record containing the prefix in 
the specified field in the query, or false if there are no matches.  After calling Seek, Next 
will return the first record greater then or equal to the prefix, and Prev will return the 
first record less then the prefix.

See also: [Query](<Query.md>),
[SeekCursor](<SeekCursor.md>)