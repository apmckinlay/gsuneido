### HorzControl

``` suneido
(control, ... [overlap:])
```

Takes a list of controls and lays them out in a horizontal row, side by side.

Xmin (width) and Xstretch are set to the total of the controls' values.

Ymin (height) and Ystretch are set to the maximum of the controls' values.

If displayed with a larger width, the excess space is distributed between the controls based on their Xstretch values.

If a control has a y stretch of 0 then it will stretch vertically to fill the Horz, but will not necessarily make the Horz itself stretch vertically.

If **overlap:** is true, then the controls are overlapped by one pixel. This eliminates the "double" line between fields.

Derived from [Group](<Group>).

See also: [VertControl](<VertControl.md>)