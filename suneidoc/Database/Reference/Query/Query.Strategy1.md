<div style="float:right"><span class="builtin">Builtin</span></div>

#### Query.Strategy1

``` suneido
(@args) => string
```

Returns the strategy for the query, showing index usage. It gives the strategy that will be used for [Query1](<../Query1.md>), [QueryEmpty?](<../QueryEmpty?.md>), [QueryFirst](<../QueryFirst.md>), or [QueryLast](<../QueryLast.md>).

If the query is a table name plus named arguments, it will be treated like the fast version of Query1 or QueryExists? The possible approaches, in order of preference:
`no select: table^key(a)`
: If there is no select the first index is used.

`key(a,b) with a: ?, b: ?`
: 

`=> key: table^(a,b)`
: A key with a select for all its fields so only one record. This includes an empty key().

`index(a,b) with a: ?, b: ?`
: 

`=> just index: table^(a,b)`
: An index that uses all of the select.

`index(a,b) with a: ?`
: 

`only: table^(a,b)`
: Only one index that uses any of the select.

`key(a,b) index(c,d) with a: ?, c: ?`
: 

`multiple indexes: (a,b) (c,d)`
: Multiple indexes that use any of the select. The indexes are read in parallel.

This streamlined optimization does not look at the data, it works strictly from the available indexes and the specified select.

If the query (first) argument is not just a table name then it will be parsed and optimized like a regular query and will give the same kind of strategy result as [QueryStrategy](<../QueryStrategy.md>). However, it will optimize for reading a single record which may, for example, avoid temporary indexes. This is affected by, and may give different results for, different data.