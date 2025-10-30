### [Hwnd](<../Hwnd.md>) - Methods
`AddHwnd(hwnd)`
: Adds the hwnd to .WndProc.HwndMap. This is normally done automatically by Hwnd.CreateWindow.

`CreateWindow(className, windowName, style, exStyle=0, x=0, y=0, w=0, h=0, id=0)`
: Calls CreateWindowEx, assigns the result to .Hwnd, and enters it in .WndProc.HwndMap.

`FixedFont()`
: Sets the current font to the ANSI_FIXED_FONT stock object.

`GetFont() => hfont`
: Returns the handle for the current font, or 0 (zero) if no font has been specified.

`Mapcmd(command_name) => command_number`
: Calls 
[.Window.Mapcmd](<../Window/window.Mapcmd.md>)

`PostMessage(msg, wParam = 0, lParam = 0) => result`
: Equivalent to:  
`PostMessage(.Hwnd, msg, wParam, lParam)`

`Repaint(erase = true)`
: Specifies that the control needs to be redrawn on the screen. The actual redraw will occur when the next WM_PAINT message is sent.   
Equivalent to:  
`InvalidateRect(.Hwnd, NULL, erase)`  
See also: MSDN documentation for InvalidateRect

`SendMessage(msg, wParam = 0, lParam = 0) => result`
: Equivalent to:  
`SendMessage(.Hwnd, msg, wParam, lParam)`

`SetFont(font = "", size = "", weight = "", text = "", underline = false)`
: If *font, size, weight*, and *underline* all have default values, the font will be set to Suneido.logfont   
*size* is in points or can be relative by passing a string starting with + or -, e.g. "+3".   
*weight* is in the range 0 to 1000. 400 is normal, 700 is bold. 0 or "" will use the default weight.   
If the font cannot be constructed, SetFont will use the default gui font stock object.   
If *text* is supplied, .Xmin, .Ymin, and .Top are set based on it.

`Update() => result`
: Causes any pending screen updates to be done.   
Equivalent to:  
`UpdateWindow(.Hwnd)`  
See also: MSDN documentation for UpdateWindow