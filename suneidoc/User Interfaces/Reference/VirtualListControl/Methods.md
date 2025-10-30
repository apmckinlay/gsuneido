### [VirtualListControl](<../VirtualListControl.md>) - Methods
ClearHighlight()
: Clears all highlight color

ClearSelect()
: Clears all selections

GetColums() => object
: Returns an object containing all columns in the list.

GetLoadedData() => object
: Returns an object containing all data loaded into memory.

GetSelectedRecord() => record
: Returns the selected row, or false if a single row is not selected.

GetSelectedRecords() => object
: Returns a list of the selected rows, or false if no row is selected.

GetField(field) => value
: Returns the value of the specified field of the selected row.   
**Note:** When you call this method, you must ensure there is a single row currently selected, otherwise an exception will be thrown.

GetQuery() => string
: Returns the current query.

HighlightValues(member, values, color)
: Highlights any row where the specified member of the row object contains a value in the specified values.

Save()
: Force saving the current record. (The standalone list saves automatically.)

Seek(field, prefix)
: Jumps to the first record containing the prefix in the specified field in the query

SetColumns(columns)
: Sets the specified columns

SetField(field, value, invalidate = false)
: Sets the specified field of the current row to the value.   
**Note:** When you call this method, you must ensure there is a single row currently selected, otherwise an exception will be thrown.

SetQuery(query, columns = false, filters = false)
: Reloads the records in the list from the specified query.