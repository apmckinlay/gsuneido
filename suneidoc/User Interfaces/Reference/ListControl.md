<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/ListControl/Methods">Methods</a></span></div>

### ListControl

``` suneido
(columns = false, data = false, defWidth = false, noShading = false, 
    noDragDrop = false, highlightColor = 0x00FF9957, noHeaderButtons = false)
```

Creates a "list" control, displays **columns** and their **data**, and allows data to be added, modified and deleted.

Xstretch and Ystretch default to 1.

The **column** argument is a list of column names.  The **data** argument is an object containing nested objects (rows).  The **defWidth** argument is the default widths of the columns. If defWidth is not specified the column widths default to 100 pixels wide.

The **noDragDrop** argument determines whether or not columns in the list can be rearranged.

The **noShading** argument determines whether or not every second line in the list will be shaded.

**highlightColor** specifies the default color for highlighted rows in the list.

The **noHeaderButtons** argument determines whether or not header buttons can be clicked.

For example:

``` suneido
Window(Object('List' columns: Object('Amount', 'Unit'), 
    data: Object(
        Object(Amount: 20, Unit: 'lbs'),
        Object(Amount: 25, Unit: 'litres'), 
        Object(Amount: 15, Unit: 'kms')), 
    defWidth: false))
```

Would display something like:

![](<../../res/listcontrol_shade.gif>)

and

``` suneido
Window(Object('List' columns: Object('Amount', 'Unit'), 
    data: Object(
        Object(Amount: 20, Unit: 'lbs'),
        Object(Amount: 25, Unit: 'litres'), 
        Object(Amount: 15, Unit: 'kms')), 
    defWidth: false, noShading:))
```

Would display something like:

![](<../../res/listcontrol.gif>)