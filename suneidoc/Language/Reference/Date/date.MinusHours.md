#### date.MinusHours

``` suneido
(date) => number
```

Returns the number of minutes between two dates.

Simply does [date.MinusSeconds](<date.MinusSeconds.md>) / 3600 so it may return fractional amounts.

For example:

``` suneido
#20000303.1030.MinusHours(#20000303.0915) => 1.25
```

**Warning**: beware of doing this across daylight savings changes.


See also:
[date.MinusDays](<date.MinusDays.md>),
[date.MinusMinutes](<date.MinusMinutes.md>),
[date.MinusMonths](<date.MinusMonths.md>),
[date.MinusSeconds](<date.MinusSeconds.md>)
