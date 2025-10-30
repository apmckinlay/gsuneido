<div style="float:right"><span class="builtin">Builtin</span></div>

### QueryAlt

``` suneido
(query [, field: value ...]) => list of record
```

QueryAlt executes the query in a very simple way. It will throw an exception if the results are too large. It does not handle schema tables or summarize. It also does not handle all uses of rules.

Named arguments (field: value) add a where to the query. For example:

``` suneido
('tables', tablename: value)
```

is equivalent to:

``` suneido
('tables where tablename = ' $ Display(value))
```

**Note**: QueryAlt should only be used for testing and fuzzing. It is **much** slower than the real query implementation.