### Writing a Hwnd Control

[Hwnd](<../Reference/Hwnd.md>) controls "wrap" a Windows control for use in Suneido user interfaces. For example, here is a simplified version of [ButtonControl](<../Reference/ButtonControl.md>):

**`MyButtonControl`**
``` suneido
Hwnd
    {
    Name: "Button"
    Xmin: 50
    Ymin: 25
    New(name = "Button")
        {
        .CreateWindow("button", name,
            WS.VISIBLE | WS.TABSTOP | BS.PUSHBUTTON)
        .Map = Object()
        .Map[BN.CLICKED] = "CLICKED"
        }
    CLICKED()
        {
        .Send("Button_Clicked")
        }
    }
```

Note:

-	inherit from Hwnd
-	specify a default name
-	create the Win32 control
-	map any notifications that you want to handle
-	handle the notifications, possibly by Send'ing messages to its controller


Note: Hwnd controls cannot handle WM window messages. If this is required, use [WndProc](<../Reference/WndProc.md>) instead.

See also:
[Writing a WndProc Control](<Writing a WndProc Control.md>)