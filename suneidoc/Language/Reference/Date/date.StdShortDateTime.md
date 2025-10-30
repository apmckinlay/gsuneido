#### date.StdShortDateTime

``` suneido
() => string
```

Returns the date formatted as yyyy-MM-dd HH:mm.

Unlike [date.ShortDateTime](<date.ShortDateTime.md>), this is not affected by the Windows settings

For example:

``` suneido
Date().StdShortDateTime()
    => "2005-08-21 11:27"
```