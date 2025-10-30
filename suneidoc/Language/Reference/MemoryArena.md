<div style="float:right"><span class="builtin">Builtin</span></div>

### MemoryArena

``` suneido
( ) => number
```

Returns the current size of the heap memory *arena* in bytes. At any given time, some of this arena will be *free* i.e. available for re-use.

For example:

``` suneido
MemoryArena()
    => 1156344
```

**Note:** The arena size normally will not shrink even if memory usage decreases.

Use [Suneido.GoMetric](<Suneido/Suneido.GoMetric.md>) to get the amount of memory currently used:

``` suneido
Suneido.GoMetric("/memory/classes/heap/objects:bytes")
```