## Size and Stretch

If a control has a known fixed size you can simply specify Xmin: and Ymin: in the control.

``` suneido
Control
    {
    Xmin: 100
    Ymin: 100
```

If the size is fixed but is initially calculated, you can set Xmin and Ymin in New.

``` suneido
Control
    {
    New(...)
        {
        .Xmin = ...
        .Ymin = ...
```

If the size changes dynamically, then you need a Recalc method that sets Xmin and Ymin. Recalc is called bottom-up (children before parent) by Window Refresh

``` suneido
Control
    {
    Recalc()
        {
        .Xmin = ...
        .Ymin = ...
        }
```

Xstretch and Ystretch can be set similarly. They are usually known fixed values (commonly either 0 or 1) except in the case of controls that contain other controls (containers) in which case the stretch usually comes from the children. It will still normally be fixed, except in the case where children are added/removed/modified dynamically.