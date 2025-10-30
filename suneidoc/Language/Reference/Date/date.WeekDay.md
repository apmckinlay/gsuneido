<div style="float:right"><span class="builtin">Builtin</span></div>

#### date.WeekDay

``` suneido
(firstDay = 'Sun') => number
```

Returns the day of the week in the range 0 to 6, where 0 is Sunday.

Optionally, you can specify which day is the *first* day of the week, either by name (upper or lower case, prefixes accepted) or by number (0 to 6).

For example:

``` suneido
#20030331.WeekDay() => 1

#20030331.WeekDay('mon') => 0
```