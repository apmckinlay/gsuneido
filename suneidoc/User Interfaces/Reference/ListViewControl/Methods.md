### [ListViewControl](<../ListViewControl.md>) - Methods
`AddColumn(name)`
: 

`AddItem(label, image = 0, lParam = 0)`
: 

`Addrow(object)`
: 

`CheckAll(checked = true)`
: This method should only be used if the listview has the LVS_EX_CHECKBOXES style. Sets the check state of all items.

`DeleteAll()`
: 

`DeleteAllColumns()`
: Deletes all of the listview columns.

`DeleteItem(item = false)`
: Deletes the item at the specified index, or the selected item if item is false.

`Disable()`
: Disables the listview window.

`Enable()`
: Enables the listview window.

`EnableMenu(enabled = true)`
: Enables or disables the right-click context menu.

`EnsureVisible(i)`
: Ensures that the item at position i is visible in the client area of the listview.

`FindItem(field, value)`
: Given the field name (column) and the value to search for, returns the index 
of the first matching item, or false if no match is found.

`GetCheckState(i)`
: Returns the check state of the item at index i.

`GetCheckedItems()`
: Returns an object containing indexes of items whose check state is true.

`GetColumnWidth(column_name) => number`
: Returns the width of the specified column (name or index).

`GetColumnWidths() => object`
: Returns a list containing the widths of all columns in the current Listview.

`GetItemPosition(i) => point`
: 

`GetLastRow()`
: Returns the item index of the last row (item) in the list.

`GetPreviousItem()`
: Returns the item index of the previous selection.

`GetRowCount()`
: Returns the number of rows (items) in the list.

`GetSelected()`
: Returns the item index of the current selection.

`GetRecord(item)`
: If there is a model (i.e. query was specified) then .model.Getrecord is returned, otherwise .Getrow is returned. This method is useful because it returns the complete model record rather than just the columns in the list which ListView will truncate to a maximum of 256 characters.

`Getrow(item)`
: Uses the index passed to return an object containing the column values for that item.

`HideColumn(column)`
: Sets the width of the specified column to 0, which hides the column. The column argument should be a column name and not an index.

`OrderColumns(column_list)`
: Column_list is specified as an object containing column indexes 
in the order they are to appear in the listview.   
For example OrderColumns(#(2 1 0)) would cause the third column to be displayed first, followed by the second, and then the first column would be displayed as the last column.

`Reset()`
: Destroys and recreates the listview window.

`ResetQuery(query, columns = false)`
: This should only be used when a query is being used as the listview's contents (virtual).  
This method is used to change the query being displayed in the list. If the columns are different, they can be specified as an object in the columns parameter.

`ResetTimer()`
: Resets the listview's timer, or creates a new one.  This is used with the timed_destroy option that destroys the listview window after a specified amount of time.

`SelectItem(i)`
: Selects the item having index i.

`SetCheckState(i, state)`
: Sets the check state of the item whose index is i.

`SetColumnName(column, text)`
: Sets the column header text for the specified column to the value passed in text.

`SetColumns(cols)`
: 

`SetColumnValue(column, text, index = false)`
: Column is the column name, text is the new text for the cell, and index 
is the index of the item whose column value will be changing.  If no index is passed, the index will be the current selection.

`SetColumnWidth(column, width)`
: Sets the width of the specified column (name or index).

`SetExtendedStyle(extended_style)`
: Sets the extended style of the listview window.

`SetImageList(images, which)`
: 

`SetItemPosition(i, x, y)`
: 

`SetListSize(size)`
: This method is used to tell the listview how many items it contains if it is a virtual listview (has the LVS_OWNERDATA style). The listview is then able to set it's scroll properties accordingly.

`SetStyle(style)`
: 

`SetView(view)`
: 

`ToObject()`
: Creates an object from the listview contents.

`UnSelect()`
: Unselects the current selection.

`UnSelectItem(i)`
: Unselects the item at the index position i.

`UpdateCachedRow(i, rec)`
: Updates the item i in the cached records for the listview display.  Also redraws the listview contents to reflect the record change.