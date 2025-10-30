#### number.RoundToPrecision

``` suneido
(precision) => number
```

Rounds the number to the specified precision.

The precision should be an integer greater than zero.

For example:

``` suneido
12345.RoundToPrecision(2) => 12000
12345.RoundToPrecision(4) => 12350
12345.RoundToPrecision(6) => 12345

.12345.RoundToPrecision(2) => .12
.12345.RoundToPrecision(4) => .1235
.12345.RoundToPrecision(6) => .12345
```


See also:
[number.Round](<number.Round.md>),
[number.RoundDown](<number.RoundDown.md>),
[number.RoundUp](<number.RoundUp.md>),
[number.RoundToNearest](<number.RoundToNearest.md>),
[number.Ceiling](<number.Ceiling.md>),
[number.Floor](<number.Floor.md>),
[number.Int](<number.Int.md>)
