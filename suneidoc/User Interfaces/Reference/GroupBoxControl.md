### GroupBoxControl

``` suneido
(title, control)
```

Creates a layout for **control**, draws a line border around it, and adds **title** (if it is specified) on the top line.

The **title** is a string containing the text of the title, and the **control** is a Control object to add the title to.

For example: 

``` suneido
Window(#(GroupBox 'Information'
    ( Form
        ( id group: 0 ) (amount group: 1 ) nl
        ( control group: 0 ) nl
        ( field group: 0 ) ( cost group: 1 )
    )
))
```

Will produce something like:

![](<../../res/groupboxcontrol.png>)