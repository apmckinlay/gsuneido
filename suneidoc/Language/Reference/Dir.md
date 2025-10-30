<div style="float:right"><span class="builtin">Builtin</span></div>

### Dir

``` suneido
(path = "*.*", files = false, details = false, block = false) => object
```

Dir returns a list of file and sub-directory names matching a given path specification.
Sub-directory names are identified by a trailing "/".

If **files** is true, sub-directories are omitted from the list.

If **details** is true, each entry in the directory is returned as an object containing the name, the date the file was last written to, and the size in bytes. For directories, the name will end with '/'.

For example:

``` suneido
Dir("c:/suneido/*.*")

=>  #("building.txt", "installing.txt", "mybook.su", "mycontacts.su", "mylib.su", 
    "release010912.txt", "SciLexer.dll", "source/", "suneido.db", "suneido.exe", 
    "translatelanguage.su",)
```

If you specify a block it will be called for each file instead of returning a list. The block can use break and continue as in other loops. For example:

``` suneido
Dir() { Print(it) }
```

Dir may throw exceptions starting with "Dir:"

See also:
[DirExists?](<DirExists?.md>),
[FileExists?](<FileExists?.md>)