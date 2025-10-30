### YesNo

``` suneido
(prompt = "", title = "", hwnd = 0, flags = 0)
```

Calls [Alert](<Alert.md>), adding MB.YESNO to flags. Returns true if the user clicks on "Yes", false otherwise.

For example:

``` suneido
if not YesNo("Continue?")
    return
```