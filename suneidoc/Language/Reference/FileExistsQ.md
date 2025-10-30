<div style="float:right"><span class="builtin">Builtin</span></div>

### FileExists?

``` suneido
(filename) => true or false
```

Returns true if the specified filename exists, false otherwise if the file does not exist.Throws an exception for an invalid path.

For example:

``` suneido
FileExists?("c:/suneido/suneido.db")
    => true

FileExists?("//nonexistent/temp/file.txt")
   => throws FileExists?: CreateFile //nonexistent/temp/file.txt: The network name cannot be found.
```

See also:
[Dir](<Dir.md>),
[DirExists?](<DirExists?.md>)