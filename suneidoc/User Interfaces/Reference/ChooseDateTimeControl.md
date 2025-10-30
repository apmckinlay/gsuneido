### ChooseDateTimeControl

Uses [ChooseField](<ChooseField.md>) to display a field with a choose button on the right side. The button will bring up a [DateTimeControl](<DateTimeControl.md>) in a dialog. Choosing OK from the dialog will fill in the selection into the field. Cancelling from the dialog will not modify the fields contents.

For example:

``` suneido
Window(#(ChooseDateTime))
```

Would display:

![](<../../res/ChooseDateTime.png>)

and clicking on the drop down button would display:

![](<../../res/datetime.png>)