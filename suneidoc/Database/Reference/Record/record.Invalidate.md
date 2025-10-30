<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Invalidate

``` suneido
( @fieldnames )
```

Marks the fields as *invalid* - i.e. the next time it is accessed, it's rule will be evaluated. 

This will also trigger any observers.

This also invalidates any dependents and triggers any observers for them.