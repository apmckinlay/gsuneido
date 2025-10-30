### ShortestKey

``` suneido
( query_string )
( list )
( query_or_cursor )
 => string
```

Returns the key with the fewest fields. Accepts either a query string, a list of keys, an open query, or an open cursor.

For example:

``` suneido
ShortestKey("stdlib")
    => "num"
```

Note: if the query has where clause that is "unique" (i.e. identifies a single record) then you'll get an empty key.

See also: [QueryKeys](<QueryKeys.md>), 
[query.Keys](<Query/query.Keys.md>)