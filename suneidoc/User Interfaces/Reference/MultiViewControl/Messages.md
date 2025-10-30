### [MultiViewControl](<../MultiViewControl.md>) - Messages

MultiViewControl sends the following notification messages to its Controller:
`MultiView_BeforeRecord(x)`
: Whenever the record being displayed in the Access or VirtualList screen. Sends the record currently loaded in the screen.

`MultiView_ItemSelected(record)`
: When the Access or VirtualList sets its data. Sends the record being set.

`MultiView_RecordChange(member, record)`
: When the record in the Access or VirtualList is changed (any of the fields are modified).

`MultiView_AfterField(field, value)`
: When the user modifies the record in the Access or VirtualList. Sends the field that was modified, and the new value of the field.

`MultiView_BeforeSaving()`
: When the Access or VirtualList record is about to be saved, after the validation is done on the record, AND the transaction has not started. 

`MultiView_BeforeSave(t)`
: When the Access or VirtualList record is about to be saved, after the validation is done on the record. Sends the transaction to be used for the save.

`MultiView_AfterSave(t)`
: After the record has been saved. Sends the transaction that was used for the save.

`MultiView_AfterSaving`
: After the record has been saved AND the transaction has been completed.

`MultiView_BeforeDelete(t)`
: Sent just before deleting a record in Access or VirtualList. Sends the transaction to be used for the delete.