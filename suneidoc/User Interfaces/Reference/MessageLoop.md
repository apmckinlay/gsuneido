<div style="float:right"><span class="builtin">Builtin</span></div>

### MessageLoop

``` suneido
(hdlg)
```

This is Suneido's built-in message loop. The main message loop is run automatically at startup. MessageLoop is normally only used by Dialog.

**hdlg** must be supplied.

The main processing is:

``` suneido
if (hdlg && GetWindowLong(hdlg, GWL_USERDATA) == 1)
    return ;
HWND window = GetAncestor(msg.hwnd, GA_ROOT);
if (HACCEL haccel = (HACCEL) GetWindowLong(window, GWL_USERDATA))
    if (TranslateAccelerator(window, haccel, &msg))
        continue ;
if (IsDialogMessage(window, &msg))
    continue ;
TranslateMessage(&msg);
DispatchMessage(&msg);
```

To set the accelerators for the dialog use:

``` suneido
SetWindowLong(hdlg, GWL.USERDATA, haccels)
```

To end the dialog use:

``` suneido
PostMessage(.Hwnd, WM.NULL, END_MESSAGE_LOOP, END_MESSAGE_LOOP)
```

**Note**: This is handled automatically by [Dialog](<Dialog.md>)