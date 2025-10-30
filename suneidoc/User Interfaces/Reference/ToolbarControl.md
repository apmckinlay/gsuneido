### ToolbarControl

``` suneido
(button, ...)
```

Creates a row of icon buttons.

Buttons can be specified with just a string (e.g. "Copy") or as an object where the first un-named value is the name of the button. e.g. #("Copy", drop:). The object form is used to pass additional options:
drop:
: true/false. In addition to receiving the standard click event, setting this attribute to receive the "Drop_"+command event when the button is right clicked. You can handle this event and display a popup menu.

"" inserts a separator (vertical bar).

">" inserts Fill e.g. to place the following buttons at the right hand end of the tool bar. This will also make the tool bar stretchable.

Toolbar uses .Window.Commands to get the shortcut, tooltip and icon information.

The icon can be one of the pre-defined Suneido font icons (see IconFont.icons) or any .emf files in imagebook. If not specified, the default icon name is the command in lowercase.

Normally a toolbar is placed at the top of the window, for example:

``` suneido
Controls: (Vert
    (Toolbar Cut Copy Paste)
    ...
    StatusBar)
```

You also have to specify the commands:

``` suneido
Commands:
    (
    (Cut,   "Ctrl+X")
    (Copy,  "Ctrl+C")
    (Paste, "Ctrl+V")
    )
```

See: User Interfaces - [Commands](<../Commands.md>)

Example:

``` suneido
Controller
    {
    Title: 'test ToolBar'
    Commands: (
        (Cut "Ctrl+X")
        (Copy ,"Ctrl+C")
        (Paste ,"Ctrl+V")
        (Parent_Folder, "Ctrl+F", "This is a test", "folder") // "folder" is the icon name
        (vsplit))

    Controls: (Vert
        (Skip 3)
        (Horz
            (Toolbar (Cut, drop:), Copy, Paste, "", Parent_Folder, ">", vsplit)
            (Vert (EtchedLine before: 0 after: 0) Field))
        (Skip 3)
        (Statusbar))

    On_Parent_Folder()
        {
        .Vert.Status.Set('On_Parent_Folder')
        }

    Drop_Cut(@unused)
        {
        .Vert.Status.Set('Drop_Cut')
        }
    }()
```