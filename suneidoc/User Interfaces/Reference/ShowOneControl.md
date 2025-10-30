### ShowOneControl

``` suneido
(control1, control2, ...)
```

Displays one of a list of controls.

The controls must be (or contain) a single control that works with RecordControl, i.e. sends 'Data' and has a Get method.

If all the controls are empty (have a value of ""), then the first control is displayed, otherwise the <u>last</u> control that has a value is displayed. For example, if you have 3 controls and the 2nd and 3rd have a value, then the 3rd will be displayed.

See also:
[HideEmptyControl](<HideEmptyControl.md>),
[ShowHideControl](<ShowHideControl.md>)