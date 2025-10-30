<div style="float:right"><span class="builtin">Builtin</span></div>

### DirExists?

``` suneido
(dirname) => true or false
```

Returns true if the specified directory exists, false if the directory does not exist. Throws an exception for an invalid path.

For example:

``` suneido
DirExists?("c:/suneido")
    => true

DirExists?("//nonexistent/temp")
   => throws DirExists?: CreateFile //nonexistent/temp: The network name cannot be found.
```

See also:
[Dir](<Dir.md>),
[FileExists?](<FileExists?.md>)