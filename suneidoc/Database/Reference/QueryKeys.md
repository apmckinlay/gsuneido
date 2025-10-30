### QueryKeys

``` suneido
(query) => object
```

Returns an object containing the keys for the query.

Note: if the query has where clause that is "unique" (i.e. identifies a single record) then you'll get an empty key.