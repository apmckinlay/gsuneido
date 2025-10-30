## Commands

Controllers can implement commands that the user can choose in several different ways:

-	menus
-	toolbar buttons
-	accelerator keys


First, the commands must be defined:

``` suneido
Commands: (
     (Users_Manual, "F1", "Open the Suneido Help", Help)
     ...
```

For each command, you can specify:

-	command "name"
-	shortcut key (or "" if the command doesn't have one)(displayed on the menu and on the toolbar tooltips)
-	status bar explanation(displayed on the status bar when the command is selected on the menu)
-	the toolbar bitmap name (if different from the command name)


The menu is then specified using just the command names with the other information used automatically from the Commands specifications.

``` suneido
Menu:
    (
    ("&Edit",
        "Cu&t", "&Copy", "&Paste", "&Delete")
    ("&Help",
        "&Users Manual", "&About Suneido")
    )
```

Notice the use of & to mark the underlined mnemonic letter for the menu commands. Any spaces in the menu strings will automatically be translated to underscores when looking up the commands. Shortcut keys will automatically be displayed next on the menus.

Toolbars are also specified using the command names. The Command information is then used to determine the bitmap to use (the default is to use the bitmap with the same name as the command). Tooltips are automatically created consisting of the command name plus the shortcut key (if there is one).

``` suneido
Controls:
    (Vert
        (Toolbar, Cut, Copy, Paste)
        Editor)
```

The command methods are prefixed with "On_".  For example:

``` suneido
On_Users_Manual()
    {
    BookControl("suneidoc");
    }
```