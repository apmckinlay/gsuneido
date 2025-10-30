#### date.Format

``` suneido
(format) => string
```

Calls [date.FormatEn](<date.FormatEn.md>) and then translates the day and month names.

For example:

``` suneido
Date().Format('dddd MMM dd, yyyy')
    => "vendredi aoû 12, 2005"
```

**Note:** The translatelanguage table must contain translations for the month and day names to enable translation.