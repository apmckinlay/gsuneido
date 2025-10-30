<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/HeaderControl/Methods">Methods</a></span></div>

### HeaderControl

``` suneido
(@args)
```

where **args** is an array containing field names and style (optional).

Creates a Windows Header control.

For example:

``` suneido
Window(Object('Header' 'Name', 'Address', 'Information' 
    style: HDS.BUTTONS))
```

would display

![](<../../res/headerControl.gif>)