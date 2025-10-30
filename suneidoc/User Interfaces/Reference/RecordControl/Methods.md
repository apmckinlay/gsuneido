### [RecordControl](<../RecordControl.md>) - Methods
`AddObserver(function)`
: Adds the function as an observer to the data record (using
[record.Observer](<../../../Database/Reference/Record/record.Observer.md>))   
If the data record is changed (using Set()) the observers are added to the new record.   
Note: The observers are not removed from the old record.

`RemoveObserver(function)`
: Undo AddObserver

`AddSetObserver(function)`
: The specified function will be called at the end of Set

`RemoveSetObserver(function)`
: Undo AddSetObserver

`Dirty?(state = '') => boolean`
: The dirty flag is set to true when the data is modified either by the application or the controls (i.e. the user). It is set to false by SetData()

`Get() => object`
: Returns the current data object.

`GetField(name) => value`
: 

`GetControl(name)`
: Returns the actual reference to the control of the provided *name* or false if there is none.

`HasControl?() => boolean`
: 

`Set(record)`
: Sets the data to the supplied record. Updates the controls. Any control whose value is not specified in the supplied object is set to "".  Sets Dirty? to false. **Note:** A record is required - an object will not work.

`SetField(name, value)`
: Used by application to set the specified name to the specified value. Calls control.Set to update the corresponding control. Sets Dirty? to true.

`Valid(forceCheck = false) => true or string`
: Returns true if the record control is not dirty or if all the controls are valid. Otherwise it returns "Invalid: " followed by a comma separated list of field prompts or headings. To force valid checking on all the controls even when the record control is not dirty, pass the forceCheck argument as true.   
Another way to force validation is to set the record control to dirty, for example:
``` suneido
.Data.Dirty?(true)
```