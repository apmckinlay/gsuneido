### [TabControl](<../TabControl.md>) - Methods
`Count() => number`
: Returns the total number of tabs.

`GetData(i)  => object`
: Returns data of tab at position *i*.

`GetImage(i) => image`
: Returns the image for tab at position *i*.

`GetSelected() => number`
: Returns the index of the selected tab.

`GetText(i) => string`
: Returns the name of the tab at posistion *i*.

`Insert(i, text, data = #(tooltip: ""), image = -1)`
: Insert tab at position *i*.

`Remove(i)`
: Deletes tab at position *i*.

`SetData(i, data)`
: Set *data* for tab at position *i*.

`SetImageList(image)`
: Sets the image for the tab.

`SetPadding(hpadding, vpadding)`
: Set horizontal and vertical padding on tabs.

`SetText(i, text)`
: Set *text* for the tab at position *i*.