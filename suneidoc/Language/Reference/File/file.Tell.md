<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Tell

``` suneido
() => offset
```

Returns the current read/write position in the file.
Tell is normally used with
[file.Seek](<file.Seek.md>)

For example, to determine the size of a file (in bytes):

``` suneido
File("tmp")
    {|f|
    f.Seek(0, "end")
    length = f.Tell()
    }
```

Note: Tell is invalid if the file was opened with mode "a"