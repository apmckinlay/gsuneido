<div style="float:right"><span class="builtin">Builtin</span></div>

### MemoryAlloc

``` suneido
( ) => number
```

Returns the Go runtime metric "/gc/heap/allocs:bytes". This is the cumulative sum of memory allocated to the heap by the application. This is just allocation, it does not go down when memory is garbage collected.

For example:

``` suneido
ReadableSize(MemoryAlloc())
    => 77.9 gb
```

See also: [Suneido.GoMetric](<Suneido/Suneido.GoMetric.md>)