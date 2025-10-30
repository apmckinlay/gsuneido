<div style="float:right"><span class="builtin">Builtin</span></div>

#### date.Plus

``` suneido
(years:, months:, days:, hours:, minutes:, seconds:, milliseconds:) => date
```

Returns a copy of the date with the specified units added.  
To subtract, give negative amounts.

For example:

``` suneido
Date("dec 31 1999").Plus(days: 1) => #20000101
Date("jan 1 2000").Plus(days: -1) => #19991231
```