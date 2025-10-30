### [HeaderControl](<../HeaderControl.md>) - Methods
`AddItem(field, width = false, tip = false)`
: Add *field* to the end of header.

`DeleteItem(idx)`
: Removes item at postion *idx* from header.

`GetItemWidth(idx) => number`
: Returns the width of item at position *idx*.

`GetNumItems() => number`
: Returns the number of items in the header.

`InsertItem(idx, field, width = false, tip = false)`
: Inserts *field* at the position *idx*.

`SetButtons(boolean)`
: Sets style of fields to HDS.BUTTONS if true.

`SetItemWidth(idx, width)`
: Sets size of item at postion *idx* to *width*.

`SwapItems(idx1, idx2)`
: Switches items at positions *idx1* and *idx2*.