#### date.MinusMinutes

``` suneido
(date) => number
```

Returns the number of minutes between two dates.

Simply does [date.MinusSeconds](<date.MinusSeconds.md>) / 60 so it may return fractional amounts.

For example:

``` suneido
#20000303.1030.MinusMinutes(#20000303.1015) => 15
```

**Warning**: beware of doing this across daylight savings changes.


See also:
[date.MinusDays](<date.MinusDays.md>),
[date.MinusHours](<date.MinusHours.md>),
[date.MinusMonths](<date.MinusMonths.md>),
[date.MinusSeconds](<date.MinusSeconds.md>)
