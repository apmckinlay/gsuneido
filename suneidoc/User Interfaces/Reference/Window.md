### Window

``` suneido
(control, title, x = CW_USEDEFAULT, y = CW_USEDEFAULT,
    w = CW_USEDEFAULT, h = CW_USEDEFAULT, 
    style = WS.OVERLAPPEDWINDOW, exStyle = 0, 
    wndclass = "SuBtnfaceArrow", show = true, 
    exitOnClose = false, keep_placement = false)
```

Top level user interface window.  Derived from [WndProc](<WndProc.md>) (via WindowBase).

-	Uses the "SuBtnfaceArrow" window class by default. The standard window classes are registered in the Init function.
-	Sets the window title to the title argument if supplied, otherwise to the control's Title
-	Creates a menu from the control's Menu if it has one
-	Redirects commands (menu, toolbar, shortcut key) to its control.
-	Saves and restores the focus when the window is inactivated
-	Calls Activate, Inactivate, MenuSelect, and Destroy on its control.
-	Responds to WM_SIZE by calling Resize on each control.
-	if x, y, w, h are <= 1 they are treated as fractions of the screen size


If the **exitOnClose** argument is true, then closing the window will cause Suneido to exit (the process is terminated).

Specify **keep_placement** as true in order to save the window size and position when it is destroyed, and restore the saved size when it is opened. The size is stored per user in the keylistviewinfo table with "wp:" $ title as the key.

##### Menus

A menu specification consists of a list of menus. The first item in each menu is the name. An item that is an empty string ("") is taken as a separator. If a menu item is an object, it will become a sub-menu. For example:

``` suneido
Menu: ((File New Open Save Exit) (Edit Undo "" Copy Cut Paste))
```

If a sub-menu has **no** items (just the name) it is *dynamic*. When the user selectst the sub-menu, Window calls "Menu_" $ name in the control to get the contents of the sub-menu. This method should return the contents of the menu. If the user selects an item, from a dynamic sub-menu, Window calls "On_" $ name and passes it the item that was chosen. For example:

``` suneido
Window(Controller
    {
    Controls: (Center (Static hello))
    Menu: ( (Menu One (SubMenu) Two) )
    Menu_SubMenu()
        { return LibraryTables() }
    On_SubMenu(option)
        { Alert(option) }
    })
```