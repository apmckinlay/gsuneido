### StatusbarControl

Creates a Windows Common Controls "statusbar".  Automatically resizes to fit the containing window.  The control name defaults to "Status".

Methods:
`AddPanel(size)`
: add a panel box (panel box position in the statusbar is 0,1,2,3 ...)

`Set(text, at)`
: Sets the text of a panel box at the indicated position

`Get(at) &rarr; string`
: returns the text of a panel box at the indicated position

`SetBkColor(color)`
: set an RGB background color of all the status bar

`SetFont(font,size)`
: set font and size of all the panles boxes of the status bar

`RemPanel(at)`
: remove panel box at the indicated position

`GetNumPanels() &rarr; number`
: returns the number of addded panels boxes

`SetPanels(panels)`
: set an object containing the panels end position in points
(e.g. #( 110, 170, 280, 450) to define 4 panels box with 
widths of 110, 60, 110, 170 pixels)

`GetPanels() &rarr; object`
: returns an object containing the panels end positions in pixels
(e.g. #( 110, 170, 280, 450))

`GetPanelRect(panel) &rarr; rect`
: returns a Rect object containing the coordinates of the rectangle
of the statusbar part indexed by panel position

Example:

``` suneido
Controller
    {
    Title: 'test'
    Xmin: 450   // window width
    Ymin: 200   // window height
    New()
        {
        .statbar = .Vert.Status         // control path
        .statbar.AddPanel(110)          // panel box at position 0
        .statbar.AddPanel(60)           // panel box at position 1
        .statbar.AddPanel(110)          // panel box at position 2
        .statbar.AddPanel(170)          // panel box at position 3
        .statbarcolor=false
        .myChooseList = .Vert.ChooseList 
        .myChooseList.SelectItem(0)      // select first item
        .myChooseMany = .Vert.ChooseMany
        .myChooseMany.Set('1,2,3,4')
        .m = 'this field have a font size of '
        .myChooseList.SetStatus(.m$'8')  //tip
        .myChooseMany.SetStatus(.m$'8')
        .myStatic1 = .Vert.Static1
        .myStatic1.Set("#panels boxes "$.statbar.GetNumPanels()$" - "$
                    "panels end posistion in points #("$
                            .statbar.GetPanels().Join(", ")$")")
        .myStatic2 = .Vert.Static2
        .myStatic2.SetColor(RGB(0, 0, 128)) // RGB(x,y,z) return blue color
        .changeAllFonts(8)
        .statbar.Set(.myChooseList.Get(),1)
        .statbar.Set(.myChooseMany.Get(),2)
        .getpanelbox1()
        }
     Controls:
        (Vert
            (ChooseList list: #('one', 'two', 'three', 'four'))
            (Skip)
            (ChooseMany list: #('1', '2', '3', '4'),
                        listDesc: #('red', 'green', 'blue', 'grey'),
                        columnWidths: #('value':30, 'desc': 50))
            (Skip)
            (Horz
                (Button 'Set font size 10' command: 'Size10' xmin:95)
                (Skip 12)
                (Button 'Set font size 8' command: 'Size8' xmin:95)
                (Skip 110)
                (Button 'Change status bar color' command: 'changeSBC')
            )
            (Skip 40)
            (Static name: "Static1", color: 8388608, xmin: 450) // blue color
            (Skip)
            (Static name: "Static2", xmin: 450)
            (Statusbar)
        )
     On_Size10()
        {
        .changeAllFonts(10)
        }
     On_Size8()
        {
        .changeAllFonts(8)
        }
    changeAllFonts(size)
        {
        font = 'MS Sans Serif'
        .statbar.SetFont(font,size)
        .myChooseList.SetFont(font,size)
        .myChooseList.SetStatus(.m$size) // field tip
        .myChooseMany.SetFont(font,size)
        .myChooseMany.SetStatus(.m$size) // field tip
        .myStatic1.SetFont(font,size)
        .myStatic2.SetFont(font,size)
        .statbar.Set('Font size '$size, 0)
        }
    getpanelbox1()
        {
        .myStatic2.Set("the content of panel box #1 is '"$.statbar.Get(1)$"'")
        }
    On_changeSBC()
        {
        .statbarcolor = (.statbarcolor is true)? false:true
        if .statbarcolor is true
            .statbar.SetBkColor(RGB(255, 0, 0)) // red color
        else
            .statbar.SetBkColor(false)  // default color
        }
    NewValue(value, source)
        {
        // source is the the control that generate the new value
        // value is a new value
        try
            {
            .statbar.Set("last selection at "$Date().Format("HH:mm:ss"), 3)
            if source.Name is 'ChooseList'
                {
                .statbar.Set(source.Get(), 1)
                .getpanelbox1()
                }
            if source.Name is 'ChooseMany'
                .statbar.Set(source.Get(), 2)
            }
        }
    }
```