### EnsureDirectories

``` suneido
EnsureDirectories(path)
```

Creates all subdirectories in a path recursively. This is equivalent to the Unix `mkdir -p` command or Go's `os.MkdirAll` function.

Unlike [EnsureDir](<EnsureDir.htm>) which creates only a single directory, EnsureDirectories will create all necessary parent directories in the specified path.

Path can be a regular path or UNC path (starting with `//`). If the path doesn't end with `/`, the last segment is assumed to be a filename.

For example:

``` suneido
EnsureDirectories('c:/folder1/folder2/folder3')
// Creates: c:/folder1, c:/folder1/folder2

EnsureDirectories('c:/folder1/folder2/')
// Creates: c:/folder1, c:/folder1/folder2

EnsureDirectories('//server/share/folder/')
// Creates: //server/share, //server/share/folder
```


See also:
[CopyFile](<CopyFile.md>),
[CreateDir](<CreateDir.md>),
[DeleteDir](<DeleteDir.md>),
[DeleteFile](<DeleteFile.md>),
[DeleteFiles](<DeleteFiles.md>),
[EnsureDir](<EnsureDir.md>),
[MoveFile](<MoveFile.md>)
