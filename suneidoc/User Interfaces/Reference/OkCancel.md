### OkCancel

``` suneido
(prompt = "", title = "", hwnd = 0, flags = 0)
```

Calls [Alert](<Alert.md>), adding MB.OKCANCEL to flags. Returns **true** if the user clicks on OK, **false** otherwise.

For example:

``` suneido
if not OkCancel("Continue?")
    return
```