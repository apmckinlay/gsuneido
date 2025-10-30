#### string.Escape

``` suneido
() => string
```

Converts: 

-	linefeed, return, and tab to '\n', '\r', '\t'
-	other non-printable characters to hex


This is the opposite of [string.Unescape](<string.Unescape.md>)

Note: If you test this from the WorkSpace you will see double backslashes, e.g. '\\t' because Display also escapes backslashes. But the actual result string only contains a single backslash. You can verify this with string.Size()