<div style="float:right"><span class="builtin">Builtin</span></div>

#### runPiped.Read

``` suneido
(nbytes = false) => string or false
```

Returns the next nbytes, or the rest of the data if nbytes is false, or false if at the end.

**Note:** Prior to BuiltDate 20241219, the default value for nbytes was 1024