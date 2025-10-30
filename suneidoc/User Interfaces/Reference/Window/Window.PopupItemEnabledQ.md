#### Window.PopupItemEnabled?

``` suneido
(hsubmenu, itempos = 0) => true or false or -1
```
`hsubmenu`
: is a menu handle

`itempos`
: is the item position in the submenu (0,1,2 ...)

Returns true if enabled, false if not. If it fails it return -1

See also:
[Window.CheckPopupItem](<Window.CheckPopupItem.md>),
[Window.EnablePopupItem](<Window.EnablePopupItem.md>),
[Window.GetPopupItemInfo](<Window.GetPopupItemInfo.md>),
[Window.PopupItemChecked?](<Window.PopupItemChecked?.md>),
[Window.SetPopupItemInfo](<Window.SetPopupItemInfo.md>),
[window.GetPopupHandle](<window.GetPopupHandle.md>)