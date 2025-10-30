<div style="float:right"><span class="builtin">Builtin</span></div>

#### number.RoundUp

``` suneido
(number) => number
```

Returns this number rounded <u>up</u> to the specified number of decimal places.

(-x).RoundUp(d) is equivalent to -(x.RoundUp(d))

For example:

``` suneido
n = 5.4321
n.RoundUp(3) => 5.433
n.RoundUp(1) => 5.44
n.RoundUp(0) => 6

(-1.3).RoundUp(0) => -2
```

Negative numbers of decimal places will round to that power of 10. For example:

``` suneido
n = 1234
n.RoundUp(-1) => 1240
n.RoundUp(-2) => 1300
```


See also:
[number.Round](<number.Round.md>),
[number.RoundDown](<number.RoundDown.md>),
[number.RoundToNearest](<number.RoundToNearest.md>),
[number.RoundToPrecision](<number.RoundToPrecision.md>),
[number.Ceiling](<number.Ceiling.md>),
[number.Floor](<number.Floor.md>),
[number.Int](<number.Int.md>)
