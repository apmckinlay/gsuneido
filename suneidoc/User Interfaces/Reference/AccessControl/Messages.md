### [AccessControl](<../AccessControl.md>) - Messages

AccessControl sends the following notification messages to its Controller:
`Access_BeforeRecord(x)`
: Whenever the record being displayed in the Access changes (including new records). Sends the record currently loaded in the Access.

`Access_SetRecord(record)`
: When the Access sets its RecordControl's data. Sends the record being set.

`Access_NewRecord(data: data)`
: When a new record is loaded into the AccessControl (user chooses the "New" button, or On_New method is called from the application). Note: dirty is set to false after this is Sent. Note: The argument is named so you use the parameter it must be called "data".

`Access_AfterNewRecord(record)`
: When a new record is loaded into the AccessControl (user chooses the "New" button, or On_New method is called from the application).  This message is sent after the Access is completely finished initializing the new record (dirty flag has been set to false already).

`Access_AllowDelete()`
: When user chooses the "Delete" option from the "Current" button. Controller should return true to allow the delete, false to prevent the delete.

`AccessBeforeDelete(t)`
: Sent just before deleting a record in AccessControl. Sends the transaction to be used for the delete.

`Access_AllowRestore()`
: When user chooses the "Restore" option from the "Current" button. Controller should return true to allow the restore, false to prevent the restore.

`Access_Restore(new_record?)`
: When Access is about to re-read the current record from the database, discarding any changes to the record that hadn't been saved. Sends a boolean value indicating whether it was a new record or not.

`Access_RecordChange(modified_fields)`
: When the record in the Access is changed (any of the fields are modified). Sends a list of the fields that were modified.

`Access_AfterField(field, value)`
: When the user modifies the record in the Access. Sends the field that was modified, and the new value of the field.

`AccessBeforeSave(t)`
: When the Access record is about to be saved, after the validation is done on the record. Sends the transaction to be used for the save.

`AccessAfterSave(t)`
: After the record has been saved. Sends the transaction that was used for the save.

`AccessAfterSaving`
: After the record has been saved AND the transaction has been completed.

`Access_CloseWindowConfirmation`
: Sent when user is attempting to close the window containing the Access. Return true to allow the window to close, false to prevent the window from closing.