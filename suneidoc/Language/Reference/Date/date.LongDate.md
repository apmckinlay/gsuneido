#### date.LongDate

``` suneido
() => string
```

Returns the date formatted using Suneido.LongDateFormat, which is set by Init from the long date format in the Windows Control Panel - Regional Settings - Date.

For example:

``` suneido
Date().LongDate() => "Sun, Nov 26, 2000"
```