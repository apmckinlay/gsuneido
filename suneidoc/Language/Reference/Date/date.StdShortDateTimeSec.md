#### date.StdShortDateTimeSec

``` suneido
() => string
```

Returns the date formatted as yyyy-MM-dd HH:mm:ss

Unlike [date.ShortDateTimeSec](<date.ShortDateTimeSec.md>), this is not affected by the Windows settings

For example:

``` suneido
Date().StdShortDateTimeSec()
    => "2005-08-21 11:27:54"
```