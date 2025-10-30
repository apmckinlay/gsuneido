### MonthCalDatesDialog

``` suneido
(dates = "", protectBefore = "")
```

Uses the Windows [MonthCalControl](<MonthCalControl.md>) common control to display a calendar from which the user can select dates.  The dates are then displayed in the [FieldControl](<FieldControl.md>) beneath the calendar.
dates
: Used to pass in a list of dates to be initially selected.  Default is "".

protectBefore
: Used to prevent the user from selecting any date that is previous to the date passed in.  If no date is passed in, the user may select any date.  Default is "".

![](<../../res/MonthCalDates.png>)