## Using RecordControl

**Category:** User Interface

**Problem**

You want to work with a number of data controls.

**Ingredients**

[RecordControl](<../User Interfaces/Reference/RecordControl.md>), 
[Controller](<../User Interfaces/Reference/Controller.md>), 
[FieldControl](<../User Interfaces/Reference/FieldControl.md>), 
[ButtonControl](<../User Interfaces/Reference/ButtonControl.md>), 
[Inspect](<../Language/Reference/Inspect.md>), 
[Alert](<../User Interfaces/Reference/Alert.md>)

**Recipe**

Enter this in a library (e.g. mylib) as My_RecordControl:

``` suneido
Controller
    {
    Controls:
        #(Record
            (Vert
                (Field name: city)
                (Field name: country)
                Skip
                (Horz
                    (Button Canada) Skip
                    (Button Home) Skip
                    (Button City) Skip
                    (Button Inspect) Skip
                    (Button Protect))))
    On_Canada()
        {
        .Data.SetField('country', 'Canada')
        }
    On_Home()
        {
        .Data.Set(#{city: 'Saskatoon', country: 'Canada'})
        }
    On_City()
        {
        Alert(.Data.GetField('city'))
        }
    On_Inspect()
        {
        Inspect(.Data.Get())
        }
    On_Protect()
        {
        control = .Data.GetControl('city')
        control.SetReadOnly(not control.GetReadOnly())
        }
    }
```

You can then run it directly from Library View or from the WorkSpace with:

``` suneido
My_RecordControl()
```

or:

``` suneido
Window(My_RecordControl)
```

Note that in On_Home where we call set we must pass a record (using #{...} or Record(...)) rather than an object (using #(...) or Object(...)).

**Discussion**

A RecordControl is a PassthruController that tracks the contents of its controls in a record. RecordControl's default name is "Data" which is why we refer to it as ".Data". A RecordControl "mirrors" the contents of the controls it contains. Changing the values in the RecordControl (using Set or SetField) will update the contents the controls. Similarly, when the contents of the controls are changed (by the user) the values in the RecordControl will be updated (and can be retrieved with Get or GetField). You can also use GetControl to get a reference to a control by name, which is often more convenient than referencing controls by "path" like .Vert.city

RecordControl "wraps" a single top-level control so you must use something like a Vert or Horz or Form to contain multiple controls.

**See Also**

[Responding to a Button](<Responding to a Button.md>)