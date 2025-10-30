### [ListBoxControl](<../ListBoxControl.md>) - Methods
`AddItem(string, index = false)`
: The string is inserted at the specified index, or added to the end of the listbox.

`DeleteAll()`
: Empties the listbox.

`DeleteItem(index)`
: Deletes the item at position index from the listbox.

`FindString(string) => int`
: Returns the index of the item equal to or prefixed by string, or -1 if no match.

`GetCount() => int`
: Returns how many items are in the listbox.

`GetCurSel() => int`
: Returns the index of the currently selected item, or -1 if there is no selection.

`GetData(i)`
: Returns the data associated with the i'th item.

`GetSelected() => int`
: Returns the data associated with the currently selected item.

`SetColumnWidth(w)`
: Sets the width of the columns in multicolumn mode.

`SetData(i, int)`
: Sets the data associated with the i'th item to int.