<div style="float:right"><span class="builtin">Builtin</span></div>

### CreateDir

``` suneido
(dir)
```

<span class="deprecated">Deprecated</span> use [EnsureDir](<EnsureDir.md>) instead.

Create a directory.

If the directory already exists it will return "CreateDir ...: already exists", but it will **not** throw an exception if this return value is not checked.

Returns true or an error string. It is return throw i.e. if you do not use the return value and the result is not true, it will throw the error.


See also:
[CopyFile](<CopyFile.md>),
[EnsureDirectories](<EnsureDirectories.md>),
[DeleteDir](<DeleteDir.md>),
[DeleteFile](<DeleteFile.md>),
[DeleteFiles](<DeleteFiles.md>),
[EnsureDir](<EnsureDir.md>),
[MoveFile](<MoveFile.md>)
