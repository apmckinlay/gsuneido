<div style="float:right"><span class="builtin">Builtin</span></div>

### ResourceCounts

``` suneido
() => object
```

Returns an object with named members containing the current counts of different types of resources. Only non-zero counts are included. This is used for tracking down resource leaks.

This has been superseded by [Suneido.Info](<Suneido/Suneido.Info.md>) and [Suneido.GoMetric](<Suneido/Suneido.GoMetric.md>)