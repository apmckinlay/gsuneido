<div style="float:right"><span class="builtin">Builtin</span></div>

### GetTempFileName

``` suneido
(path, prefix) => string
```

Returns a temporary file name.

For example:

``` suneido
GetTempFileName(GetTempPath(), "dbg")
    => C:/Users/Andrew/AppData/Local/Temp/dbgB831.tmp
```

See also: [GetTempPath](<GetTempPath.md>)