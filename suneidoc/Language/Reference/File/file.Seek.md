<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Seek

``` suneido
(offset, origin = "set") => true or false
```

Sets the current read/write position in the file.
Returns true if successful, false otherwise.
Seek is often used along with
[file.Tell](<file.Tell.md>)

`origin` can be:
**`"set"`**
: offset is relative to the beginning of the file i.e. *absolute*.

**`"end"`**
: offset is relative to the end of the file.

**`"cur"`**
: offset is relative to the current position in the file.

Note: Seek is invalid if the file was opened with mode "a"