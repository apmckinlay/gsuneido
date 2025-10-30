<div style="float:right"><span class="stdlib">stdlib</span></div>

### AddFile

``` suneido
(filename, @strings)
```

Appends the strings to the file, creating it if necessary.

``` suneido
AddFile("tmp.txt", "hello", " ", "world", "\n")
```

Passing multiple arguments is more efficient than concatenating them, especially for long strings.

#### See Also

[GetFile](<GetFile.md>),
[PutFile](<PutFile.md>),
[File](<File.md>),