<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.Complete

``` suneido
()
```

May throw "transaction.Complete failed ..."

Note: Only update transactions (that have done updates) can fail to Complete. (Due to conflicts with other update transactions.)

See also:
[transaction.Rollback](<transaction.Rollback.md>)