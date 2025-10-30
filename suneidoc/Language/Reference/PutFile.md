### PutFile

``` suneido
(filename, @strings)
```

Writes (or overwrites) the contents of the file with the strings.

``` suneido
PutFile("tmp.txt", "hello", " ", "world", "\n")
```

Passing multiple arguments is more efficient than concatenating them, especially for long strings.

See also:
[AddFile](<AddFile.md>),
[GetFile](<GetFile.md>),
[File](<File.md>)