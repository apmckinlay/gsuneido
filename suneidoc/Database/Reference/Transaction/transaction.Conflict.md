<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.Conflict

``` suneido
() => string
```

Returns a string describing the source of a transaction conflict.

For example:

``` suneido
if not tran.Complete()
    Alert(tran.Conflict())
```