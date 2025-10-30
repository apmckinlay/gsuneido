<div style="float:right"><span class="builtin">Builtin</span></div>

### ExePath

``` suneido
() => string
```

Returns the path of the running Suneido executable.

For example:

``` suneido
ExePath()
    => "C:\\Suneido\\suneido.exe"
```

**Note:** The backslashes are not actually double - since backslashes are special characters, the representation for an actual backslash is two backslashes.