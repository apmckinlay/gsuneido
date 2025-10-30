<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.GoMetric

``` suneido
(string = false) => number
```

Returns a [Go runtime metric](<https://pkg.go.dev/runtime/metrics>) or a list of the available metrics.

Currently only handles metrics that return float64 or uint64.

These are also available from the web monitor.

See also: [MemoryAlloc](<../MemoryAlloc.md>), [Suneido.Info](<Suneido.Info.md>)