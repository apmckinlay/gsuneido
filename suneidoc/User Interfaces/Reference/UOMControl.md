<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/UOMControl/Methods">Methods</a></span></div>

### UOMControl

``` suneido
(numctrl, uomctrl, readonly = false, div = '')
```

where **numctrl** is usually a [NumberControl](<NumberControl.md>) type, and **uomctrl** is usually a [ChooseListControl](<ChooseListControl.md>) or a [KeyControl](<KeyControl.md>) type.

This control takes in two controls and lays them out in a horizontal row where **numctrl** is on the right and **uomctrl** is on the left (with div between them).

For example:

``` suneido
Window(#(UOM (Number mask: '#########.##')
    (ChooseList #('liters', 'gallons'))))
```

would display

![](<../../res/uomControl.gif>)

User can enter a value into the field on the left and select a unit of measure (uom) by clicking the down arrow button on the right.

The value of a UOMControl (e.g. what Get returns) is:

``` suneido
numctrl.Get() $ ' ' $ uomctrl.Get()
```