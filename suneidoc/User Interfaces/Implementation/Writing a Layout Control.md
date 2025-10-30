### Writing a Layout Control

For example, here is a control that will center a non-stretchable control within a stretchable area:

**`MyCenterControl`**
``` suneido
Container
    {
    Xstretch: 1
    Ystretch: 1
    New(control)
        {
        .ctrl = .Construct(control)
        .Xmin = .ctrl.Xmin
        .Ymin = .ctrl.Ymin
        }
    Resize(x, y, w, h)
        {
        xs = w - .ctrl.Xmin
        ys = h - .ctrl.Ymin
        .ctrl.Resize(x + xs / 2, y + ys / 2, .ctrl.Xmin, .ctrl.Ymin)
        }
    GetChildren()
        {
        return Object(.ctrl)
        }
    }
```

Which could be used like:

``` suneido
Window(#(MyCenter (Static Hello)))
```

Note:

-	inherit from 
	[Container](<../Reference/Container.md>)
-	construct child control(s)
-	set Xmin and Ymin
-	implement GetChildren (required by Container)