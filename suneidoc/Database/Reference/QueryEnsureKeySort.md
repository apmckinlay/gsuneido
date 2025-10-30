### QueryEnsureKeySort

``` suneido
(query) => query
```

Returns the query with the key *sort* added (if there isn't one).

If the query already has a sort but it's not a key, it will throw an exception.

For example:

``` suneido
QueryEnsureKeySort("tables")
    => "tables sort table"
```

Used by [QueryApply](<QueryApply.md>) and [QueryApplyMulti](<QueryApplyMulti.md>).