### [BrowseControl](<../BrowseControl.md>) - Methods
`AddRecord(rec, idx = false)`
: Adds the record specified by **rec** at the position **idx**.If idx is not specified the record will be added as the last record in the browse. Set **browseForceAdd** to true in **rec** to force the browse to add the record. This is useful sometimes when adding records from the application code.

`Get() => object()`
: Returns an object containing nested objects (data) in the browse.

`GetColums() => object()`
: Returns an object containing all columns in the browse.

`GetCurrentRecord() => record`
: Returns the current row, or false if a single row is not selected.

`GetField(field) => value`
: Returns the value of the specified field of the current row.   
**Note:** When you call this method, you must ensure there is a single row currently selected, otherwise an exception will be thrown.

`GetProtectField() => string`
: Returns the name of the protected field.

`GetQuery() => string`
: Returns the current query.

`GetStickyFields() => object`
: Returns an object containing names of sticky fields.

`Save()`
: Saves the current record in the browse.

`SaveColumns()`
: Saves the orders and sizes for columns in the browse by user.

`Set(value)`
: Set data in browse to value.

`SetField(field, value, idx = current_row, invalidate = false)`
: Sets the specified field of the current row to the value.   
**Note:** When you call this method without specifying an idx, you must ensure there is a single row currently selected, otherwise an exception will be thrown.

`SetQuery(query, columns = false, header_data = false, t = false,
        max_records = false, max_records_msg = '')`
: Reloads the records in the browse from the specified query. If t is a valid transaction, it will be used to read the records. Use max_records and max_records_msg to limit the number of records loaded by the browse, if the number of records read exceeds the number specified in max_records, it will stop loading at max_records and display the msg from max_records_msg to the user.