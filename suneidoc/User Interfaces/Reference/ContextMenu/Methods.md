### [ContextMenu](<../ContextMenu.md>) - Methods
`Show(hwnd, x = 0, y = 0, left = false) => number`
: Activates a context menu at the specified position (using TrackPopupMenu)
and returns the one based index of the option chosen,
or zero if no item was chosen (the menu was cancelled).   
If left is true, menu options can only be chosen with the left mouse button.
Otherwise, they can be chosen with either mouse button.   
For example:

``` suneido
ContextMenu(x, y)
    {
    i = ContextMenu(#(Cut, Copy, Paste)).Show(.Hwnd, x, y)
    if (i is 0)
        return 0
    ...
```

`ShowCall(control, x = 0, y = 0) => value`
: Activates a context menu at the specified position (using TrackPopupMenu)
and if an item is chosen, calls the corresponding method in the control.
The method names are "On_Context_" $ option.   
For example:

``` suneido
ContextMenu(x, y)
    {
    ContextMenu(#(Cut, Copy, Paste)).ShowCall(this, x, y)
    }
On_Context_Cut()
    ...
On_Context_Copy()
    ....
On_Context_Paste()
    ...
```