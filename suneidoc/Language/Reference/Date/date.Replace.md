#### date.Replace

``` suneido
(year = false, month = false, day = false, 
    hour = false, minute = false, second = false, millisecond = false) => date
```

Returns a copy of a date with the specified components replaced.

For example:

``` suneido
Date().Replace(year: 1900) => #19001125.170630545
```