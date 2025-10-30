### CheckBoxControl

``` suneido
(text = "", font = "", size = "", weight = "", readonly = false,
    set = false, tip = "")
```

Creates a Windows "button" control with the "checkbox" and "lefttext" styles.
text
: the label on the checkbox

font, size, weight
: let you override the default font style

readonly
: if true the control will initially be disabled

set
: if true the checkbox will initially be checked

tip
: a tooltip for the check box

**Note:** Control.Construct automatically sets the text
on CheckBox fields with prompts.

CheckBoxControl has Set, Get, and Dirty? methods to work with RecordControl.
It also has a Protect method to work with __protect rules.

For example:

``` suneido
Window(#(CheckBox 'Allow Bonus'))
```

would produce something like:

![](<../../res/checkbox.gif>)