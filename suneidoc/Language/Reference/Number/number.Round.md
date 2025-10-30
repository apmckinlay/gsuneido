<div style="float:right"><span class="builtin">Builtin</span></div>

#### number.Round

``` suneido
(number) => number
```

Returns this number rounded to the specified number of decimal places using "half up".

(-x).Round(d) is equivalent to -(x.Round(d))

For example:

``` suneido
n = 123.456
n.Round(2) => 123.46
n.Round(1) => 123.5
n.Round(0) => 123

1.5.Round(0) => 2

(-.33).Round(1) => -.3
(-.66).Round(1) => -.7
```

Negative numbers of decimal places will round to that power of 10. For example:

``` suneido
n = 7654
n.Round(-1) => 7650
n.Round(-2) => 7700
```


See also:
[number.RoundDown](<number.RoundDown.md>),
[number.RoundUp](<number.RoundUp.md>),
[number.RoundToNearest](<number.RoundToNearest.md>),
[number.RoundToPrecision](<number.RoundToPrecision.md>),
[number.Ceiling](<number.Ceiling.md>),
[number.Floor](<number.Floor.md>),
[number.Int](<number.Int.md>)
