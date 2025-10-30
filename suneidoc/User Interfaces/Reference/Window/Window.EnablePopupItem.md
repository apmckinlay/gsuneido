#### Window.EnablePopupItem

``` suneido
(hsubmenu, item = 0, enable = true, byposition = true) => number
```

Enable/disable a popup menu item
`hsubmenu`
: is a menu handle

`item`
: may be an item position in the submenu (0,1,2 ...) or an ID menu item identifier

`enable`
: if true enable the item, if false disable the item

`byposition`
: if true the item value is its position, if false its ID

Returns the previous state of the menu item either `MF.DISABLED`, `MF.ENABLED`, or `MF.GRAYED`, or -1 if it fails.

See also:
[Window.CheckPopupItem](<Window.CheckPopupItem.md>),
[Window.GetPopupItemInfo](<Window.GetPopupItemInfo.md>),
[Window.PopupItemChecked?](<Window.PopupItemChecked?.md>),
[Window.PopupItemEnabled?](<Window.PopupItemEnabled?.md>),
[Window.SetPopupItemInfo](<Window.SetPopupItemInfo.md>),
[window.GetPopupHandle](<window.GetPopupHandle.md>)