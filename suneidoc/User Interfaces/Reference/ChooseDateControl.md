### ChooseDateControl

``` suneido
(mandatory = false, width = 10, checkCurrent = false, checkValid = false)
```

Uses [ChooseField](<ChooseField.md>) to display a [DateControl](<DateControl.md>) with a choose button on the right side. The button will bring up a [MonthCalControl](<MonthCalControl.md>) in a dialog displaying a calendar where the user can select a date.

For example:

``` suneido
Window(#(ChooseDate))
```

Would display:

![](<../../res/ChooseDate.png>)

and clicking on the drop down button would display:

![](<../../res/monthcal.png>)