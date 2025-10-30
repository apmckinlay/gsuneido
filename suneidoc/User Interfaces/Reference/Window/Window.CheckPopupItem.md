#### Window.CheckPopupItem

``` suneido
(hsubmenu, item = 0, check=true, byposition = true) => number
```

Check/uncheck a popup menu item.
`hsubmenu`
: is a submenu handler

`item`
: may be an item position in the submenu(0,1,2 ...) or a ID menu item identifier

`check`
: if true check the item, if false uncheck the item

`byposition`
: if true the item value is its position, if false is its ID

Returns the previous state of the menu item: `MF.CHECKED`, or `MF.UNCHECKED`, or -1 if it fails.

See also:
[Window.EnablePopupItem](<Window.EnablePopupItem.md>),
[Window.GetPopupItemInfo](<Window.GetPopupItemInfo.md>),
[Window.PopupItemChecked?](<Window.PopupItemChecked?.md>),
[Window.PopupItemEnabled?](<Window.PopupItemEnabled?.md>),
[Window.SetPopupItemInfo](<Window.SetPopupItemInfo.md>),
[window.GetPopupHandle](<window.GetPopupHandle.md>)