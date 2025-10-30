### Containment Requirements

In order to be placed in the standard containers, controls must support the following: (Most controls should support containment.)

#### Members

<div class="table-full-width">

|  |  | 
| :---- | :---- |
| Xmin | the minimum width of the control | 
| Xstretch | how much the width of the control can be stretched (0 for no stretch) | 
| Ymin | the minimum height of the control | 
| Ystretch | how much the height of the control can be stretched (0 for no stretch) | 


The stretch values of the controls in a container are summed 
and any excess room in the container is distributed between the controls proportionally.

#### Methods

``` suneido
Resize(x, y, w, h)
Update()
Destroy()
```

Note: The base Control class contains default definitions.