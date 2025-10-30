<div style="float:right"><span class="builtin">Builtin</span></div>

#### query.Keys

``` suneido
() => object
```

Returns a list of the keys for the query.

Note: if the query has where clause that is "unique" (i.e. identifies a single record) then you'll get an empty key.