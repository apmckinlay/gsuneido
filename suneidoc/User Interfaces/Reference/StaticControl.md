### StaticControl

``` suneido
(text = "", font = "", size = "", weight = "", justify = "LEFT",
    underline = false, color = "", whitebgnd = false, 
    status = "", sunken = false, notify = false, bgndcolor = "")
```

Creates a Windows "static" control.

For example, #(Static "Hello World") would produce:

![](<../../res/Static.gif>)

#### Parameters
`text`
: is the text to display, it is translated

`font`
: is the name of the font to use. If you want to specify a different font size or weight, you must also specify the font (name).

`size`
: is a numeric font size

`weight`
: a number from 0 to 1000 where 400 is normal and 700 is bold

`justify`
: "LEFT" or "CENTER" or "RIGHT" justify the text on xmin value

`underline`
: true/false, underline or not the displayed text

`color`
: may be a number e.g. color: 0xff0000 (blue) or an object that contain 3 numeric RGB values e.g. color: #(255, 0, 0) (red)

`whitebgnd`
: true/false, set the background color to white

`status`
: allows to set a text to display as a tip when the control has the focus

`sunken`
: is true to design a sunken border to the text

`notify`
: is true to trap mouse click and double click over the control

`bgndcolor`
: Sets the color of the background. May be either a number or an object with RGB values (like color).

#### Sends
`Static_Click`
: If notify is enabled and the user clicks on the text.

`Static_DoubleClick`
: If notify is enabled and the user double clicks on the text.

#### Methods
`Get()`
: returns the original untranslated text.

`Set(text)`
: changes the text dynamically, is translated

`SetFont(font = "", size = "", weight = "", underline = false)`
: to change font, size, weight or underline

`SetColor(color)`
: Set the text color to a specified value, either a number or an object containing 3 RGB values

`SetBgndBrush(flag)`
: Set the background to one of the stock object brushes e.g. SO.WHITE_BRUSH, SO.LTGRAY_BRUSH, SO.GRAY_BRUSH, SO.DKGRAY_BRUSH, SO.BLACK_BRUSH

`SetBgndColor(color)`
: Set the background to a specified color, either a number or an object containing 3 RGB values

`SetEnabled(enabled)`
: true/false to set enabled/desabled the control

`SetVisible(visible)`
: true/false to set visible/invisible the control

`Update()`
: update the control

`Repaint(erase = true)`
: repaint the control

#### Examples

Note: Remove the final () if putting this code in a library.

``` suneido
Controller
    {
    Title: 'test StaticControl'
    Xmin: 200
    Ymin: 300
    New()
        {
        .Vert.St1.SetColor(RGB(0, 0, 255))
        .Vert.Hz1.St4.Set('Hello World')
        .Vert.Hz1.St4.SetFont('Ms Sans Serif', 20, underline:)

        .Vert.Hz2.StOne.SetStatus('click to change back color')
        .Vert.Hz2.StTwo.SetStatus('doubleclick to change text color')
        .oneclick = -1
        .twoclick = -1
        }
    Commands:
        ((Exit, "Ctrl+X"))
    Controls:
        (Vert
            (Skip)
            (Horz
                Fill
                (Static 'Hello World', justify: 'CENTER', 
                    font: 'Ms Sans Serif', size: 14, weight: 700)
                Fill
            )
            (Horz
                Fill
                (Static 'Hello World', justify: 'CENTER',
                    font: 'Ms Sans Serif', size: 14, whitebgnd: true)
                Fill
            )
            (Skip)
            (Static name: 'St1', 'Hello World', xmin: 100)
            (Field width: 11)
            (Static 'Hello World', color: #(0, 0, 255), justify: "CENTER", xmin: 100)
            (Field width:11)
            (Static 'Hello World', color: 16711680, justify: "RIGHT", xmin: 100)
            (Field width: 11)
            (Skip 20)
            (Horz name: "Hz1"
                (Static name: 'St4', color: 16711680, justify: 'CENTER')
                (Skip 20)
                (MenuButton "change to", #("invisible", "visible", "disabled", "enabled"))
            )
            (Skip 30)
            (Horz name: "Hz2"
                (Static name: 'StOne', 'click', font: 'Arial', size: 14, justify: 'CENTER'
                        sunken:, notify:, whitebgnd:, xmin: 200)
                (Skip)
                (Static name: 'StTwo', 'double click', font: 'Arial', size: 14, justify: 'CENTER'
                        sunken:, notify:, xmin: 200)
            )
        )

    On_change_to_invisible()
        { .Vert.Hz1.St4.SetVisible(false) }

    On_change_to_visible()
        { .Vert.Hz1.St4.SetVisible(true) }

    On_change_to_desabled()
        { .Vert.Hz1.St4.SetEnabled(false) }

    On_change_to_enabled()
        { .Vert.Hz1.St4.SetEnabled(true) }

    Static_Click(source)
        {
        if source.Name is 'StOne'
            {
            brushes = Object(SO.LTGRAY_BRUSH, SO.DKGRAY_BRUSH, SO.GRAY_BRUSH,
                              SO.BLACK_BRUSH, SO.WHITE_BRUSH)
            ++.oneclick
            if (.oneclick is 5)
                .oneclick = 0
            if (.oneclick is 4)
                .Vert.Hz2.StOne.SetColor(RGB(0, 0, 0))
            else
                .Vert.Hz2.StOne.SetColor(RGB(255, 255, 0))
            .Vert.Hz2.StOne.SetBgndBrush(brushes[.oneclick])
            }
        }
     Static_DoubleClick(source)
        {
        if source.Name is 'StTwo'
            {
            colors = Object(255, 32768, 16711680, 33023, 32896, 8388863, 0)
            ++.twoclick
            if (.twoclick is 7)
                .twoclick = 0
            .Vert.Hz2.StTwo.SetColor(colors[.twoclick])
            }
        }
    }()
```

Another example, demonstrating colors:

``` suneido
Controller
    {
    Xmin:300
    Ymin:300
    Title: 'test StaticControl with colors'
    New()
        {
         .oldcolor= CLR.RED
         .Vert.S1.SetColor(CLR.WHITE)
         .Vert.Horz.S5.SetBgndColor(CLR.YELLOW)
         .Vert.Horz.S5.SetColor(RGB(0, 0, 255))
         .Vert.Horz.S7.SetBgndBrush(SO.LTGRAY_BRUSH)
        }
    Commands:
        ((Exit,"Ctrl+X"))
    Controls:
        (Vert
            Skip
            (Static 'Suneido' name: 'S0' size:12)
            Skip
            (Static 'Suneido' name: 'S1' bgndcolor:#(0, 0, 0) size:12)
            Skip
            (Static 'Suneido, double click here!' name: 'S2' color: 3107669 weight:800
                       xmin:250 notify:true sunken:true size:12)
            Skip
            (Static 'Suneido' name: 'S3' bgndcolor:#(222, 184, 135) size:12)
            Skip
            (Static 'Suneido' name: 'S4' color:0x003C14DC size:12)
            Skip
            (Horz
            (Static 'Suneido' name: 'S5' xmin:100 size:12)
            (Static 'Suneido' name: 'S6' whitebgnd:true size:12)
            (Static 'Suneido' name: 'S7' size:12)
            )
        )
    Startup()
        {.Window.Center()}
    Static_DoubleClick(source)
        {
         if source.Name is 'S2'
            {
              color = .Vert.S2.GetColor()
             .Vert.S2.SetColor(.oldcolor)
             .oldcolor = color
            }
        }
    }()
```