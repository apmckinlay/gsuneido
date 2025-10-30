<div style="float:right"><span class="builtin">Builtin</span></div>

#### query.RuleColumns

``` suneido
() => object
```

Returns a list of the rule (unsaved) columns in the query (uncapitalized).

**Warning**: For complex queries this is not exact. For example, with union and join it will only include columns that are rules on both sides.

See also: 
[query.Columns](<query.Columns.md>),
[query.Keys](<query.Keys.md>)