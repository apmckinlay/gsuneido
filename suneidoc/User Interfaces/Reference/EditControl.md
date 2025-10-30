### EditControl

``` suneido
( mandatory = false, readonly = false, style = 0 )
```

Creates a Windows "edit" control.  The style argument is added to the style argument to CreateWindow.

If **mandatory** is true, the EditControl can not be left empty by the user.

If **readonly** is true, the field appears greyed out and the user can not enter any data into the field.

EditControl supports Get and Set for DataControl and Dirty? for ExplorerControl.  It also handles cut, copy, paste, and undo commands.  Below is a list of its methods.

``` suneido
    On_Cut()
    On_Copy()
    On_Paste()
    On_Undo()
    EN_KILLFOCUS()
    Get()
    Set(value)
    Valid?() 
    SelectAll()
    Dirty?(dirty = "")
    SetReadOnly(readOnly)
    GetReadOnly()
    Destroy()
```

See also: 
[FieldControl](<FieldControl.md>), 
[EditorControl](<EditorControl.md>)