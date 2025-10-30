<div style="float:right"><span class="builtin">Builtin</span></div>

#### date.MinusSeconds

``` suneido
(date) => number
```

Returns the number of seconds between two dates with millisecond accuracy. e.g. 1.5 for one second plus 500 milliseconds.

**Warning**: beware of doing this across daylight savings changes.

For example:

``` suneido
Date("10:00").MinusSeconds(Date("9:59")) => 60
```


See also:
[date.MinusDays](<date.MinusDays.md>),
[date.MinusHours](<date.MinusHours.md>),
[date.MinusMinutes](<date.MinusMinutes.md>),
[date.MinusMonths](<date.MinusMonths.md>)
