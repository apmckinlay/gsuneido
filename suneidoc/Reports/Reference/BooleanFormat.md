### BooleanFormat

``` suneido
(data = "", w = false, width = false, justify = "left",
    font = false, yesno = false, strings = #("true", "false"),
    color = false)
```

Prints a boolean (true/false) value. Normally, "true" and "false" are printed. If **yesno** is true, then "yes" and "no" will be used. Or you can specify your own **strings**.

Derived from [TextFormat](<TextFormat.md>)

See also:
[CheckBoxFormat](<CheckBoxFormat.md>),
[CheckMarkFormat](<CheckMarkFormat.md>)