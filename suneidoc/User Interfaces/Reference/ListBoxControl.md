<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/ListBoxControl/Methods">Methods</a></span></div>

### ListBoxControl

``` suneido
(string [, string ...], multicolumn = false )
```

Creates a Windows "listbox" control and calls AddItem for each of the supplied strings.

Xstretch and Ystretch default to 1.

If multicolumn is true, the ListBox will have a horizontal scroll
bar and multiple columns of items.

For example:

``` suneido
Window(#(ListBox one two three))
```

![](<../../res/List.gif>)

Send's:
`ListBoxSelect(i)`
: Sent when LBN_SELCHANGE is received.

`ListBoxDoubleClick(i)`
: Sent when LBN_DBLCLK is received.

`ListBox_ContextMenu(x, y)`
: Sent when WM_CONTEXTMENU is received. x and y are in screen coordinates. If from a mouse right click, the item under the mouse (if there is one) will be selected.