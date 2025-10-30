### SpinnerControl

``` suneido
(rangefrom = 0, rangeto = 99999, width = false, set = false,
    mask = "##,###", justify: "RIGHT", status = "", mandatory = false,
    euro = false, increase = 1, rollover = false)
```

Creates a NumberControl with UpDown arrows on the side for changing the value.
`euro`
: if true, a EuroNumberControl is used instead of a NumberControl

`increase`
: the amount that the up and down arrows will change the value

`rollover`
: if true, the up and down arrows will "wrap around", i.e. up from rangto will change the value to rangefrom, down from rangefrom will change the value to rangeto

![](<../../res/Spinner.gif>)

SpinnerControl passes its arguments on to 
[NumberControl](<NumberControl.md>).