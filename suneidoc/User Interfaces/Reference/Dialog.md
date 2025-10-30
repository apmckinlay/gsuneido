<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/Dialog/Methods">Methods</a></span></div>

### Dialog

``` suneido
(parentHwnd, control, style = 0, exstyle = 0, border = 5, title = '', 
    x = '', y = '', posRect = false, keep_size = false) => result
```

Derived from [WndProc](<WndProc.md>) via [WindowBase](<WindowBase.md>). Similar to [Window](<Window.md>), but uses [MessageLoop](<MessageLoop.md>) to create a modal dialog. All other windows are disabled while the dialog is open. This is not always convenient but it prevents errors related to nested message loops.

The control must call .Window.Result(result) to end the dialog and return a result.

The ESCAPE key will normally call On_Cancel. There is a default definition in WindowBase:

``` suneido
On_Cancel()
    {
    .Window.Result(false)
    }
```

If you have a Cancel button then it will call the same method. The Escape key is normally equivalent to a Cancel button.

The dialog is centered over the parent window.

Only certain combinations of style and exstyle work. For example:

<div class="table-style table-full-width">

|  |  | 
| :---- | :---- |
| close button on the title bar (default (if no style is specified) | <code>style: WS.SYSMENU \| WS.CAPTION<br />exStyle: WS_EX.DLGMODALFRAME</code> | 
| no title bar | <code>style: WS.POPUP \| WS.BORDER // for 2d border<br />style: WS.POPUP \| WS.MODALFRAME // for 3d border</code> | 
| thin title bar with close button, resizable, double click on title bar maximizes | <code>style: WS.DLGFRAME \| WS.SYSMENU \| WS.SIZEBOX \| WS.MAXIMIZEBOX,<br />exStyle: WS_EX.TOOLWINDOW</code> | 

</div>

If a dialog is closed with the close button on the title bar, it returns false.

For example:

**`MyDialog`**
``` suneido
Controller
    {
    Title: "Test Dialog"
    Controls:
        (Vert
            name
            Skip
            OkCancel
            )
    On_OK()
        { .Window.Result(.Vert.Pair.name.Get()) }
    }
```

OK is normally the default button, triggered by the ENTER key. You can specify a different default button as in:

``` suneido
Controls: (... (Button 'Go') ...)
DefaultButton: Go
```

**Note**: Do <u>not</u> use the defaultButton option in ButtonControl - this only affects appearance, not behavior.

Which can then be used like:

``` suneido
result = Dialog(.Window.Hwnd, MyDialog)
```

NULL (0) can be supplied for the parent hwnd if the dialog is not controlled by another window.

The **x** and **y** parameters can be used to specify the location of the top left corner of the dialog box on the screen.  If they are not specified the dialog box will be displayed in the center of the parent window, or if the parentHwnd is 0, in the center of the screen.

If **posRect** is specified, the dialog is positioned above, to the right, to the left, or below, depending on where it will fit within the working area of the screen.

Specify **keep_size** as a string in order to save the dialog size when it is destroyed, and restore the saved size when it is opened. The size is stored per user in the keylistviewinfo table with the specified string as the key.

A [Controller](<Controller.md>) that is designed to be used in a Dialog commonly defines a CallClass like:

``` suneido
Controller
    {
    CallClass(hwnd = 0)
        {
        Dialog(hwnd, this)
        }
```