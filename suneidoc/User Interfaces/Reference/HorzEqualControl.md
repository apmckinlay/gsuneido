### HorzEqualControl

``` suneido
(@controls, pad = 20)
```

Makes all the ButtonControl's the same width as the widest buton. Adds the specified padding to each button's width.

For example:

``` suneido
Window(#(HorzEqual (Static eg) Skip (Button short) (Button reallyreallylong)))
```

![](<../../res/horzequal.png>)

See also: [HorzEvenControl](<HorzEvenControl.md>)