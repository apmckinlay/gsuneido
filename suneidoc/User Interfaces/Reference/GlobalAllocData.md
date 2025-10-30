<div style="float:right"><span class="builtin">Windows Builtin</span></div>

### GlobalAllocData

``` suneido
(string) => handle
```

Returns a handle to a GlobalAlloc (GMEM_MOVEABLE) containing the string.

The string may contain binary data, including zero bytes

You must GlobalFree the returned handle when you are finished with it.

See also: [GlobalData](<GlobalData.md>), [GlobalAllocString](<GlobalAllocString.md>)