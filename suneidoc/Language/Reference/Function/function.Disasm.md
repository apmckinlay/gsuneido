<div style="float:right"><span class="builtin">Builtin</span></div>

#### function.Disasm

``` suneido
(source = false) => string
```

Returns a listing of the internal byte code for the function.

The listing will be byte code only unless the source is supplied. The source must match what the function was compiled from.