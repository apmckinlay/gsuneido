### VertControl

``` suneido
(control, ... [overlap:])
```

Takes a list of controls and lays them out in a vertical column, one below the other.

Ymin (height) and Ystretch are set to the total of the controls' values.

Xmin (width) and Xstretch are set to the maximum of the controls' values.

If displayed with a larger height, the excess space is distributed between the controls based on their Ystretch values.

If a control has an x stretch of 0 (instead of false) then it will stretch horizontally to fill the Vert, but will not necessarily make the Vert itself stretch horizontally. For example, to make a set of buttons all the same width:

``` suneido
(Vert
    (Button One xstretch: 0)
    (Button Two xstretch: 0)
    (Button Three xstretch: 0))
```

If **overlap:** is true, then the controls are overlapped by one pixel. This eliminates the "double" line between fields.

Derived from [Group](<Group>).

See also: [HorzControl](<HorzControl.md>)