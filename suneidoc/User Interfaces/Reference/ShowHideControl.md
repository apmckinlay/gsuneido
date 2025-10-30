### ShowHideControl

``` suneido
(control, showFunction)
```

ShowHideControl can be used to wrap a control (specified as the **control** argument) that is only to be shown under certain conditions. The **showFunction** argument must be a callable, or the global name of a callable. The callable should return true if the control is to be shown. Any other return value will cause the control not to be shown.

For example:

``` suneido
ShowHideControl(#(ImageControl, "c:/images/face.gif"),
    function () { return false }) // will not show

ShowHideControl(#(ImageControl, "c:/images/face.gif"),
    function () { return true }) // will show

// using the name of a global function (you must define this in a library)
ShowHideControl(#(ImageControl, "c:/images/face.gif"),
    "Test_ShowImage?")
```

The above examples simply use functions that return true/false. You can, of course, get as complicated as you wish with the callable that you are passing to this control

For development, you can make all ShowHide controls visible by setting:

``` suneido
Suneido.ShowAll = true
```

(and restore normal behaviour by setting it to false or deleting the member)

You can also control this from the Toggle ShowHide Show All option on the IDE menu.

See also:
[ShowOneControl](<ShowOneControl.md>)