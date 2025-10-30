### BorderControl

``` suneido
(control, border)
```

Creates a layout for **control** and adds **border** around it.

The **control** is a Control object, and the **border** is the border width.  If **border** is not specified, then the border width defaults to 10 pixels wide.

For example:

``` suneido
Window(
    Controller
        {
        Xstretch: 0
        Ystretch: 0
        New()
            {  super(.layout()) }
        layout()
            { return #(Vert
                (Border 
                    ( Form
                        ( id group: 0 ) (amount group: 1 ) nl
                        ( control group: 0 ) nl
                        ( field group: 0 ) ( cost group: 1 )
                    ), border: 20
                ))
            }
        }
    )
```

Would display something like:
![](<../../res/bordercontrol.png>)