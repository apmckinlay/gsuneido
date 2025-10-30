#### date.ShortDate

``` suneido
() => string
```

Returns the date formatted using Suneido.ShortDateFormat, which is set by Init from the short date format in the Windows Control Panel - Regional Settings - Date.

For example:

``` suneido
Date().ShortDate() => "00/11/26"
```