<div style="float:right"><span class="builtin">Builtin</span></div>

### FileSize

``` suneido
(filename)
```

Returns the size of a file.

Prior to 20241218 it returned false if the file did not exist.
As of BuiltDate 20241218 it throws "FileSize \<filename>: does not exist"

May throw an exception if there are problems accessing the file.

See also: [file.Size](<File/file.Size.md>)