### Writing a WndProc Control

[WndProc](<../Reference/WndProc.md>) controls "wrap" a Windows control for use in Suneido user interfaces. WndProc inherits from Hwnd, adding the ability to respond to WM window messages.

For example, here is a simplified version of the EtchedLine WndProc control:

**`MyEtchedLineControl`**
``` suneido
WndProc
    {
    Ymin: 6
    Xstretch: 1
    New()
        {
        .CreateWindow("SuneidoWindow", "", WS.VISIBLE)
        .SubClass()
        .shadowpen = CreatePen(PS.SOLID, 0,
            GetSysColor(COLOR.BTNSHADOW))
        .highlightpen = CreatePen(PS.SOLID, 0,
            GetSysColor(COLOR.BTNHIGHLIGHT))
        }
    PAINT()
        {
        hdc = BeginPaint(.Hwnd, ps = Object())
        GetClientRect(.Hwnd, r = Object())
        oldPen = SelectObject(hdc, .shadowpen)
        MoveTo(hdc, r.left, (r.top + r.bottom) / 2 - .5)
        LineTo(hdc, r.right, (r.top + r.bottom) / 2 - .5)
        SelectObject(hdc, .highlightpen)
        MoveTo(hdc, r.left, (r.top + r.bottom) / 2 + .5)
        LineTo(hdc, r.right, (r.top + r.bottom) / 2 + .5)
        SelectObject(hdc, oldPen)
        EndPaint(.Hwnd, ps)
        return 0
        }
    DESTROY()
        {
        DeleteObject(.shadowpen)
        DeleteObject(.highlightpen)
        return 0
        }
    }
```

Note:

-	inherit from WndProc
-	create the Win32 control
-	SubClass in order to handle WM_ windows messages
-	handle the necessary window messages e.g. define PAINT to handle WM_PAINT
-	release any resources in DESTROY


If you don't need to handle any WM window messages, use [Hwnd](<../Reference/Hwnd.md>) instead.

See also:
[Writing a Hwnd Control](<Writing a Hwnd Control.md>)