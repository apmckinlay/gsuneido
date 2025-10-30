#### window.Commands

``` suneido
( ) => object
```

Returns an object containing information about the commands.  This information is taken from the Commands member of the Window's control.  For example:

``` suneido
Controller
    {
    Commands:(
        (Users_Manual, "F1", "Open the Suneido Help", Help)
          ...
```

would become:

``` suneido
Window.Commands() =>
    Object(
        Suneido_Help: Object(
            accel: "F1",
            help: "Open the Suneido Help",
            bitmap: "Help"
            id: 123)
```

Where the id is the numeric id assigned by Mapcmd.

This information is used by Window and also by the ToolbarControl.

See also:
[User Interfaces - Commands](<../../Commands.md>)