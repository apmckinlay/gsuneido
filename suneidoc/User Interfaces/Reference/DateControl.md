### DateControl

``` suneido
(width = 10, status = "A date e.g. May 29 or May 29, 2003 or 5/29/03",
    readonly = false, showTime = false, mandatory = false, style = 0, checkCurrent = false,
    checkValid = false)
```

Derived from [FieldControl](<FieldControl.md>) but with validation, formatting, and conversion to [Date](<../../Language/Reference/Date/Date.md>)'s.

width, status, readonly, and style are passed directly to FieldControl. If showTime is true, then the time will be displayed as well (using [date.Time](<../../Language/Reference/Date/date.Time.md>)). If checkCurrent is true, a warning will be displayed if a date is entered that is more than 6 months in the past or more than 6 months in the future. This is useful for dates on transactions.

Dates are displayed using [date.ShortDate](<../../Language/Reference/Date/date.ShortDate.md>). However, if the resulting string would not be converted back to the same date, then the year is widened to four digits. For example, if your date setting was MM/dd/yy, Jan 1, 2005 would be displayed as 05/01/01, but Jan 1, 1922 would be displayed as 1922/01/01