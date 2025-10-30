<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.PreSet

``` suneido
(field, value)
```

Set the specifield field to the specified value, bypassing rules and observers. For example:

``` suneido
record.PreSet('age', 23)
```

is similar to:

``` suneido
record.age = 23
```

**Note:** This method should only be used in very specific circumstances.