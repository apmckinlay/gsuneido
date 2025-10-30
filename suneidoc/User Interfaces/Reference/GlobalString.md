### GlobalString

``` suneido
(handle) => string
```

Returns the string from a GlobalAlloc, excluding the terminating zero byte.

If the data does not have a terminating zero byte, the entire data will be returned.

See also: [GlobalData](<GlobalData.md>), [GlobalAllocString](<GlobalAllocString.md>)