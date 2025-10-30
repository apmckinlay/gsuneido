#### Window.SetPopupItemInfo

``` suneido
(hmenu, item = 0, byposition = true, fMask = false,
    fType = 0, fState = 0, wID = 0, hSubMenu = 0, hbmpChecked = 0,
    hbmpUnchecked = 0, dwItemData = 0, dwTypeData = "", cch = 0)
```

Returns a `MENUITEMINFO` object, or 0 if it fails.
`hmenu`
: is generally a submenu handle but may be a menu handle

`item`
: may be an item position in the submenu (0,1,2 ...) or the ID of the item

`byposition`
: if true the item value is its position, if false its ID

`fMask`
: specifies which information to set, the default is to return all the information: `MIIM.STATE | MIIM.ID | MIIM.SUBMENU | MIIM.CHECKMARKS | MIIM.TYPE | MIIM.DATA | MIIM.STRING | MIIM.BITMAP | MIIM.FTYPE`

**Note:** You have to set fMask with the exact flags matching the values you want to set.

See also:
[Window.CheckPopupItem](<Window.CheckPopupItem.md>),
[Window.EnablePopupItem](<Window.EnablePopupItem.md>),
[Window.GetPopupItemInfo](<Window.GetPopupItemInfo.md>),
[Window.PopupItemChecked?](<Window.PopupItemChecked?.md>),
[Window.PopupItemEnabled?](<Window.PopupItemEnabled?.md>),
[window.GetPopupHandle](<window.GetPopupHandle.md>)

For more information see Microsoft MSDN