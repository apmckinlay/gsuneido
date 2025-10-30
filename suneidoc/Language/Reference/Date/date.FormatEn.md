<div style="float:right"><span class="builtin">Builtin</span></div>

#### date.FormatEn

``` suneido
(format) => string
```

Converts the Date to an English string using the supplied format string in which the following characters can be used:

``` suneido
y       year
M       month
d       day
h       hour in 12 hour format
H       hour in 24 hour format
m       minute
s       second
a       am/pm
A       AM/PM
t       AM/PM
```

For names, 4 or more letters means to use the full name, 3 letters means the three letter abbreviation.  For numbers, the number of letters determines the minimum number of digits, except that yy is taken as the last two digits of the year.

For year, only 'yy' and 'yyyy' should be used. The behavior of 'y' and 'yyy' are undefined and may differ on different versions of Suneido.

Other characters are simply copied to the output string.

For example:

``` suneido
dddd, MMMM d, yyyy   Monday, February 21, 2000
ddd, MMM. dd, 'yy    Mon, Feb. 21, '00
MM/dd/yy             02/21/00
yyyyMMdd             20000221
h:mmaa               1:34pm
HH:mm:ss             13:34:09
yyyy-MM-dd H:mm      2000-02-21 13:34
```

Characters can be "escaped" with backslash or with single quotes. For example:

``` suneido
Date().Format("dd \de MMM")
    => "30 de Nov"
```
or:
``` suneido
Date().Format("dd 'de' MMM")
    => "30 de Nov"
```

See also:
[date.Format](<date.Format.md>)