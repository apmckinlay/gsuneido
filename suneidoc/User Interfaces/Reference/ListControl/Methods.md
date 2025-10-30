### [ListControl](<../ListControl.md>) - Methods
`AddHighlight(rowno, color = false)`
: Highlights the specified row. The default highlight color will be used if no color is specified.

`AddRow(row)`
: Adds a single object row to the end of the list.

`AddRows(rows)`
: Adds each row member of the object rows to the end of the list.

`ClearHighLight(rowno = false)`
: Removes all highlighting from the list. If rowno is specified, then the highlighting is removed for that row only.

`DeleteAll()`
: Removes all rows (items) from the list.

`DeleteSelection()`
: Deletes the currently selected row(s).

`Get() => object`
: Returns a reference to the actual data in the list.

`GetCol(idx) => string`
: Returns the name of column at index idx.

`GetHighlighted() => object`
: Returns the highlighted rows.

`GetNumCols() => number`
: Returns the number of columns in the list.

`GetNumRows() => number`
: Returns the number of rows (items) in the list.

`GetRow(rowno) => object`
: Returns the specified row.

`GetSelection() => object`
: Returns the currently selected rows.

`HighlightValues(member, values, color = false, sortHighlight = false, group = false)`
: Highlights any row where the specified member of the row object contains a value in the specified values. Optional a specific color can be used. The sortHighlight option will sort the list so that highlighted rows come first, and the group option will group the highlighted rows by highlight color. The group option can only be used if the sortHighlight is used.

`InsertRow(rowno, row)`
: Inserts a row at the specified rowno position.

`RowHighlighted?(rowno) => boolean`
: Returns true if the row is highlighted, false otherwise.

`Set(rows)`
: Set the list's data to rows.  Note that rows is a list of objects.

`SetRow(rowno, row)`
: Replaces the row at position rowno with the new value.

`SortHighlight(group = false)`
: Sorts the list so that highlighted rows are displayed first. If the group option isnt false, then the highlighted rows will also be grouped by highlight color.