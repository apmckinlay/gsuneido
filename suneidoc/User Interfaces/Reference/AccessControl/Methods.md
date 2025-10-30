### [AccessControl](<../AccessControl.md>) - Methods
`AccessGoto(field, value)`
: Locates and displays the record containing the specified value in the specified field.

`AccessObserver(fn, at = false)`
: Calls the specified function when setdata, save, delete or restore is done on the access. Use the at option to control the order the observers get called if there will be multiple observer functions on the same access.

`ChangeQuery(query)`
: Changes the query for the access by starting a cursor for the new query and locating the first record.  This process also starts a new transaction for the access.

`CurrentRecord() => record`
: Returns the current record from the access.

`Dirty?(dirty = '') => true or false`
: Returns true if anything has been modified in the record being displayed in the access, or false is nothing has changed. If true is passed to the method, then the record will be considered dirty (modified). If false is passed to the method, the record will not be considered dirty (not modified).

`EditMode?() => true or false`
: Returns true if the access can be edited, or false if the record is readonly (in view only mode). Some or all of the fields may still be protected even when in edit mode if the protectField option or __protect rules are used.

`GetControl(field) => control`
: Returns the control for the specified field name by calling GetControl on the Access's RecordControl.

`GetData() => record`
: Returns the current record in the access by calling Get on the Access's RecordControl.

`GetExcludeSelectFields() => object`
: returns an object containing a list of the fields that are excluded from the select's field list (users can not use these fields to select on). These fields are specified by the excludeSelectFields option for AccessControl.

`GetFields() => object`
: Returns an object containing a list of fields from the query.

`GetKeys() => object`
: Returns an object containing a list of keys from the query.

`GetOriginal() => record`
: Returns the original record read from the query, before any changes were made to it.

`GetQuery() => string`
: Returns the query used by the access.

`GetRecordControl() => RecordControl instance`
: Returns the RecordControl being used by the access.

`Get_Select_vals() => object`
: Returns an object containing the current select criteria being used by the AccessControl.

`Locate?()`
: Used to specify whether or not the locate utility is available in the access or not.

`NewRecord?() => true or false`
: returns true if the record is new (wasn't read from the query), false otherwise.

`On_Current_Delete()`
: Deletes the current record. To display a message to the user to explain why a delete is not allowed, define a global function named <Tablename>_allow_delete (for example, if the table is called equipment, the function name will be Equipment_allow_delete). If the delete is not allowed, return a message explaining the reason, otherwise return "".

`On_Current_Print()`
: Brings up a dialog to print the current record.

`On_Current_Reason_Protected()`
: If the current record is protected, and the rule that protected it also has specified the reason, it will be displayed.

`On_Current_Restore()`
: Discards any unsaved changes to the current record. Re-reads the record from the query and displays the record.

`On_Current_Save()`
: Saves the current record.

`On_Edit()`
: Switches the access between edit and view modes. In view mode the user can not modify any of the access fields.

`On_First()`
: Loads the first record from the access's query and displays it.

`On_Global_Crosstable()`
: Starts the Crosstable tool for generating a report with fields summarized as rows and columns.

`On_Global_Export()`
: Displays a dialog that allows the user to export records from the acccess's query to a file.

`On_Global_Import()`
: Displays a dialog that allows the user to import records from a file into the access's query.

`On_Global_Reporter()`
: Brings up the Reporter tool for generating custom reports.

`On_Global_Summarize()`
: Brings up a dialog to apply a summary to the query's data.

`On_Go()`
: Locates and displays the record containing the value in the field from the locate control (usually at bottom right of Access)

`On_Last()`
: Loads the last record from the access's query and displays it.

`On_Next()`
: Loads the next record from the record currently being displayed and displays it.

`On_Prev()`
: Loads the previous record from the record currently being displayed and displays it.

`On_Select()`
: Brings up a select dialog so the user can specify a selection on the query.

`RefreshQuery()`
: Rereads the current record from the query.

`Reload() => true or false`
: Rereads the current record in a new transaction and displays it. Returns true if the record was successfully reloaded, false if not.

`RemoveAccessObserver(fn)`
: Removes the function from the list of observers for the Access.

`Save() => true or false`
: Saves the current record. Returns true if the save succeeded, false if it didn't.

`SaveFor(action) => true or false`
: Used for displaying a message to the user when a particular action can not be done on a new non-dirty record. Use this when saving the access before a particular operation. Returns true if the save was successful, false otherwise.

`SetReadOnly(readOnly)`
: Pass true to disallow any editing in the access, false to allow editing. Any protectField or __protect rules will still protect  the appropriate fields even when SetReadOnly(false) is used.

`SetSelectVals(select_vals)`
: Used to set the select criteria for the Access. AccessControl saves the select criteria used last time the Access was used and applies it the next time. This is the method that is used to set the select values from the last time the access was used. This method does not actually apply the restrictions to the query.

`SetWhere(where, quiet = false) => true or false`
: Applies the specified where clause to the query. If quiet is true, a message will be displayed to the user if no records are found. Returns true if records are found, false otherwise.

`Status(status)`
: Displays the specified status in the status bar.

`Valid?() => true or false`
: Checks all of the data controls in the access to see if they are valid.  Also checks the validation rule for the access (specified by the validField option) and displays the message if the rule returned one. Returns false if anything was invalid, true otherwise.