<div style="float:right"><span class="builtin">Windows Builtin</span></div>

### GlobalAllocString

``` suneido
(string) => handle
```

Returns a handle to a GlobalAlloc (GMEM_MOVEABLE) containing the string plus a terminating zero byte.

If the string contains zero bytes, only the portion before the first will be saved.

You must GlobalFree the returned handle when you are finished with it.

See also: [GlobalString](<GlobalString.md>), [GlobalAllocData](<GlobalAllocData.md>)