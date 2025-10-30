<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/MultiViewControl/Messages">Messages</a></span></div>

### MultiViewControl

``` suneido
New(.args, .accessArgs, .listArgs, .accessGoTo? = false, defaultMode = 0)
```

Creates a control used to navigate and modify data either in form or list screen, user can use the Flip button on the top left to switch

The control is a combination of [AccessControl](<AccessControl.md>) and [VirtualListControl](<VirtualListControl.md>)

The **args** are the common arguments shared between AccessControl and VirtualListControl

The **accessArgs** are used to initiate the inside AccessControl, please refer to the [AccessControl](<AccessControl.md>) documentation for detailed information

The **listArgs** are used to initiate the inside VirtualListControl, please refer to the [VirtualListControl](<VirtualListControl.md>) documentation for detailed information

The control remembers the view mode last time used, the default view mode is form sceen, you can specify **defaultMode** as 1 to be using list mode

If the screen is access from AccessGoTo, the screen is forced to be form mode, and it does not modify the last view mode for the current screen.

For example:

``` suneido
Database('ensure sales (id, name, desc, amount, cost) key (id)')

Window(#('MultiView',
    args: #(
        'sales',
        title: 'Sales',
        option: 'Sales',
        startLast:,
        protectField: 'sales_protect',
        excludeSelectFields: #(cost)),
    accessArgs: #(
        1: #(Vert, #(Form
            (id, group: 0), (name, group: 1) nl, 
            (desc, group: 0) (amount, group: 1))
            #Customizable), 
        locate: #(keys: (id) columns: (id, name))),
    listArgs: #(
        columns: #(id, name, desc, amount),
        headerSelectPrompt:,
        enableExpand:,
        defaultColumns: #(id, name),
        filtersOnTop:)
    )
)
```

See also:
[AccessControl](<AccessControl.md>), [VirtualListControl](<VirtualListControl.md>)