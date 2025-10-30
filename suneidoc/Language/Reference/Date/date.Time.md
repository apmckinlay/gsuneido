#### date.Time

``` suneido
() => string
```

Returns the time formatted using Suneido.ShortDateFormat, which is set by Init from the time format in the Windows Control Panel - Regional Settings - Time.

For example:

``` suneido
Date().Time() => "14:21"
```