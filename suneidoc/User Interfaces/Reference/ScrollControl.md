<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/ScrollControl/Methods">Methods</a></span></div>

#### ScrollControl

``` suneido
( control )
```

Creates a scrollable window.  Scroll bars will be displayed or not depending on whether the control is larger than the window.  i.e. resizing the window to  larger than the control will cause the scroll bars to disappear.

Xstretch and Ystretch default to 1.

For example:

``` suneido
Window(#(Scroll
    (Vert
        firstname
        address
        zip_postal
        phone
        )
    ))
```