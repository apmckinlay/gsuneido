<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.GetDefault

``` suneido
(member, default_value)
```

If the object contains **member**, its value will be returned, otherwise the supplied **default_value** will be returned. For example:

``` suneido
#(a: 12, b: 34).GetDefault("b", 99) => 34

#(a: 12, b: 34).GetDefault("c", 99) => 99
```

If default_value is a block, then the result of evaluating it will be returned. This is useful if you only want the default_value expression evaluated if necessary e.g. because it is slow, or because it may fail if the member exists. For example:

``` suneido
object.GetDefault("result", { CalculateResult() })
```

See also: [object.GetInit](<object.GetInit.md>)