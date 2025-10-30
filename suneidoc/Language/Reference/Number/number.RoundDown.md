<div style="float:right"><span class="builtin">Builtin</span></div>

#### number.RoundDown

``` suneido
(number) => number
```

Returns this number rounded down to the specified number of decimal places.

(-x).RoundDown(d) is equivalent to -(x.RoundDown(d))

For example:

``` suneido
n = 123.98673
n.RoundDown(3) => 123.986
n.RoundDown(1) => 123.9
n.RoundDown(0) => 123
```

Negative numbers of decimal places will round to that power of 10. For example:

``` suneido
n = 7654
n.RoundDown(-1) => 7650
n.RoundDown(-2) => 7600
```


See also:
[number.Round](<number.Round.md>),
[number.RoundUp](<number.RoundUp.md>),
[number.RoundToNearest](<number.RoundToNearest.md>),
[number.RoundToPrecision](<number.RoundToPrecision.md>),
[number.Ceiling](<number.Ceiling.md>),
[number.Floor](<number.Floor.md>),
[number.Int](<number.Int.md>)
