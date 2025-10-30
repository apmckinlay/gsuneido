### NumberControl

``` suneido
(mask = "-###,###,###", readonly = false, rangefrom = false, rangeto = false,
    width = false, set = false)
```

Derived from
[FieldControl](<FieldControl.md>)
.
mask
: Used to format the number for display.  Default is -###,###,### 

readonly
: If true the control will initially be disabled.

rangefrom, rangeto
: The valid range of values

width
: Specifies the width of the control in digits.

set
: An initial value for the control.

**Note:** The width of the control is the maximum of the mask size and the width parameter.

For example:

``` suneido
Window(#(Number mask '####', rangefrom: 1900, rangeto: 2010))
```