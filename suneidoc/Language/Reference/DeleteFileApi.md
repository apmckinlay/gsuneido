<div style="float:right"><span class="builtin">Builtin</span></div>

### DeleteFileApi

``` suneido
(filename) => true or error
```

Delete a file. Normally it's better to use [DeleteFile](<DeleteFile.md>) which retries.

**Note**: On Windows, if the delete fails with "access denied" it will try to remove the read-only attribute (by setting the attributes to "normal") and then retry the delete. In the rare case that the second delete fails, this could leave the file existing but with the attributes changed.

If the file does not exist, it returns "DeleteFileApi: \<filename> does not exist" 
(Note: the position of the colon changed as of BuiltDate 20241218)

Returns true or an error string. It is return throw i.e. if you do not use the return value and the result is not true, it will throw the error.


See also:
[CopyFile](<CopyFile.md>),
[CreateDir](<CreateDir.md>),
[EnsureDirectories](<EnsureDirectories.md>),
[DeleteDir](<DeleteDir.md>),
[DeleteFile](<DeleteFile.md>),
[DeleteFiles](<DeleteFiles.md>),
[EnsureDir](<EnsureDir.md>),
[MoveFile](<MoveFile.md>)
