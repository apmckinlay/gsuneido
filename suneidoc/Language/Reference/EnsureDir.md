<div style="float:right"><span class="builtin">Builtin</span></div>

### EnsureDir

``` suneido
(dirname)
```

Creates a directory if it doesn't already exist.

If dirname exists but is not a directory it will return an error of "EnsureDir: {dirname} exists but is not a directory"

Returns true or an error string. It is return throw i.e. if you do not use the return value and the result is not true, it will throw the error.


See also:
[CopyFile](<CopyFile.md>),
[CreateDir](<CreateDir.md>),
[EnsureDirectories](<EnsureDirectories.md>),
[DeleteDir](<DeleteDir.md>),
[DeleteFile](<DeleteFile.md>),
[DeleteFiles](<DeleteFiles.md>),
[MoveFile](<MoveFile.md>)
