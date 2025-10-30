<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/AccessControl/Messages">Messages</a><a href="/suneidoc/User Interfaces/Reference/TabsControl/Methods">Methods</a></span></div>

### TabsControl

``` suneido
(control [, control ... ] bottom = false, themed = false,
    constructAll = false, destroyOnSwitch = false)
```

A controller for the Windows Tab control.

The controls must have a Tab: member that specifies the label to go on its tab.

For example:

``` suneido
TabsControl(#(Static 'this is one'  Tab: 'One') #(Editor, Tab: 'Editor'))
```

The bottom and themed options are passed on to [TabControl](<TabControl.md>).

By default, tab contents are constructed "lazily" i.e. when the tab is selected. If constructAll is true, then all the tab contents will be constructed right away.

If destroyOnSwitch is true, then tab contents will be destroyed when the user switches to another tab, otherwise they are not.

Example of TabsControl that pass values between the tabs:
(Each tab has an icon and can access a popupmenu)

``` suneido
Controller
    {
    Title: 'tab example'
    Xmin: 100
    Ymin: 100
    New()
        {
        .inittab = true
        .images = CreateImageList(16, 16, IDI.SUNEIDO, IDI.DOCUMENT, IDI.PROCESS)
        //
        .tabs = .Vert.Tabs
        .tabs.SetImageList(.images)
        //
        .tab0 = .tabs.WndPane  //tab 0 handler
        .f1 = .FindControl('F1')
        .st11 = .FindControl('ST11')
        .st12 = .FindControl('ST12')
        .color = 0x00ff0000
        .st11.SetColor(.color)
        .st12.SetColor(.color)
        //
        .tab1_firstime =true
        .tab1 = false //tab 1 handler, set by Startup()
        .tab2_firstime =true
        .tab2 = false //tab 2 handler, set by Startup()
        .stbar = .Vert.Status
        .stbar.AddPanel(100)
        .stbar.AddPanel(100)
        .stbar.AddPanel(200)
        }
    Startup()
        {
        //force load tab handles
        .tabs.Select(1)
        .tabs.Select(2)
        .tabs.Select(0)
        .tabs.SetImage(0, 0)
        .tabs.SetImage(1, 1)
        .tabs.SetImage(2, 2)
        .f2.Set('field 2')
        .f3.Set('field 3')
        .setstatics()
        }
    Commands:
        ((Exit,"Ctrl+X"))
     Controls:
        (Vert
            (Tabs
                (Tab: "one", Vert
                    (Skip)
                    (Static 'edit field 1')
                    (Field name: 'F1', set: 'field 1')
                    (Skip 10)
                    (Static 'field content in tab 1')
                    (Static '' name: 'ST11')
                    (Skip)
                    (Static 'field content in tab 2')
                    (Static '' name: 'ST12')
                )
                (Tab: "two", Vert
                    (Skip)
                    (Static 'edit field 2')
                    (Field name:'F2')
                    (Skip 10)
                    (Static 'field content in tab 0')
                    (Static '' name: 'ST21')
                    (Skip)
                    (Static 'field content in tab 2')
                    (Static '' name: 'ST22')
                )
                (Tab: "three", Vert
                    (Skip)
                    (Static 'edit field 3')
                    (Field name: 'F3')
                    (Skip 10)
                    (Static 'field content in tab 0')
                    (Static '' name: 'ST31')
                    (Skip)
                    (Static 'field content in tab 1')
                    (Static '' name: 'ST32')
                )
            )
            (Statusbar)
        )
    TabsControl_SelectTab(source)
        {
        if (.tabs.GetSelected() is 1 and .tab1_firstime)
            {
            .tab1_firstime = false
            .f2 = .FindControl('F2')
            .st21 = .FindControl('ST21')
            .st22 = .FindControl('ST22')
            .st21.SetColor(.color)
            .st22.SetColor(.color)
            }
        if (.tabs.GetSelected() is 2 and .tab2_firstime)
            {
            .tab2_firstime = false
            .f3 = .FindControl('F3')
            .st31 = .FindControl('ST31')
            .st32 = .FindControl('ST32')
            .st31.SetColor(.color)
            .st32.SetColor(.color)
            }
        if (not .inittab)
            .setstatics()
        else
            .inittab= false
        }
    TabClick(i)
        {
        .stbar.Set('cliked tab ' $ i, 1)
        }
    SelectTab(i)
        {
        .stbar.Set('selected tab ' $ i, 0)
        }
    setstatics()
        {
        .st11.Set(.f2.Get())
        .st12.Set(.f3.Get())
        .st21.Set(.f1.Get())
        .st22.Set(.f3.Get())
        .st31.Set(.f1.Get())
        .st32.Set(.f2.Get())
        }
    TabContextMenu(x, y)
        {
        menu = #((name: 'four'), (name: 'five'), '', (name: 'six'))
        ContextMenu(menu).ShowCall(this, x, y)
        }
    On_Context_four()
        {
        .stbar.Set('clicked menu item four', 2)
        }
    On_Context_five()
        {
        .stbar.Set('clicked menu item five', 2)
        }
    On_Context_six()
        {
        .stbar.Set('clicked menu item six', 2)
        }
    }
```