#### date.StdShortDate

``` suneido
() => string
```

Returns the date formatted as yyyy-MM-dd.

Unlike [date.ShortDate](<date.ShortDate.md>), this is not affected by the Windows settings

For example:

``` suneido
Date().StdShortDate()
    => "2005-08-21"
```