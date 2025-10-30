#### Window.GetPopupItemInfo

``` suneido
(hmenu, item = 0, byposition = true, fMask = false)
    => object or false
```

Retrieve information about a menu item.
`hmenu`
: is generally a submenu handle but may be a menu handle

`item`
: may be an item position in the submenu(0,1,2 ...) or the ID of the  item

`byposition`
: if true the item value is its position, if false its ID

`fMask`
: specifies which information to retrieve, the default is to return all the information: `MIIM.STATE | MIIM.ID | MIIM.SUBMENU | MIIM.CHECKMARKS | MIIM.TYPE | MIIM.DATA | MIIM.STRING | MIIM.BITMAP | MIIM.FTYPE`

Returns a `MENUITEMINFO` object, or 0 if it fails. The members of the returned object are:

``` suneido
#(hSubMenu: , fMask: , wID: , hbmpUnchecked: , fState: , cbSize: ,
    cch: , dwItemData: , fType: , hbmpChecked: , dwTypeData: )
```

See also:
[Window.CheckPopupItem](<Window.CheckPopupItem.md>),
[Window.EnablePopupItem](<Window.EnablePopupItem.md>),
[Window.PopupItemChecked?](<Window.PopupItemChecked?.md>),
[Window.PopupItemEnabled?](<Window.PopupItemEnabled?.md>),
[Window.SetPopupItemInfo](<Window.SetPopupItemInfo.md>),
[window.GetPopupHandle](<window.GetPopupHandle.md>)

For more information see Microsoft MSDN