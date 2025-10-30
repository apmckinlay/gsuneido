<div style="float:right"><span class="builtin">Builtin</span></div>

#### Date

``` suneido
() => current date
(date) => date
(string, format = "yMd") => date or false
(year:, month:, day:, hour:, minute:, second:, millisecond:) => date or false
```

Date represents a date and time.  Returns false if the date cannot be interpreted.

With no arguments it returns the current date and time.

If passed a string it attempts to interpret it as a date and/or time. Most reasonable formats are handled. For example:

``` suneido
Feb 14 2000
Monday 14 Feb 2000
Feb 14 '00
Feb 14
July 14
0/2/14
60/4/25
4/25/60
7-8-9 3:4:5am
3:04pm 7-8-9
20000303
20000303.1030
20000303.103000
```

Ambiguous dates such as 7/8/9 are interpreted based on the format argument,
for example:

``` suneido
Date("7/8/9", "Mdy") => #20090708)
```

Also:

-	If just a date is supplied the time is set to 00:00:00.
-	If a date is supplied without a year, then the closest year will be used.  For example, if it is January, then December would refer to the previous year.
-	If a two digit year is given, it is interpreted as within the range of 80 years back and 20 years ahead.  For example, if it is 2000, then 00 to 20 would be 2000 to 2020 and 21 to 99 would be 1921 to 1999.
-	If just a time is supplied the date is set to the current date.
-	Capitalization is not significant.
-	Day names (e.g. Tuesday) are ignored.
-	Abbreviations are interpreted as the first matching name.  For example, Ma would be March rather than May.
-	Unrecognized words or numbers will cause the interpretation to fail.
-	Punctuation is ignored except that a single quote introduces a two digit year, and colons mark times.
-	Handles dates from the year 1600.


**Note:** Date's are immutable like numbers and strings.
This means they cannot be modified once they are created.