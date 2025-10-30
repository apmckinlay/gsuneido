<div style="float:right"><span class="builtin">Builtin</span></div>

### DeleteDir

``` suneido
(dir) => true or error
```

Delete a directory and its contents including sub-directories.

If the directory does not exist it will return true.

Returns true or an error string. It is return throw i.e. if you do not use the return value and the result is not true, it will throw the error.


See also:
[CopyFile](<CopyFile.md>),
[CreateDir](<CreateDir.md>),
[EnsureDirectories](<EnsureDirectories.md>),
[DeleteFile](<DeleteFile.md>),
[DeleteFiles](<DeleteFiles.md>),
[EnsureDir](<EnsureDir.md>),
[MoveFile](<MoveFile.md>)
