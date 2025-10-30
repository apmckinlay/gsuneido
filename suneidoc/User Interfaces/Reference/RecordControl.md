<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/RecordControl/Methods">Methods</a></span></div>

### RecordControl

``` suneido
(control)
```

A PassthruController that tracks the contents of its controls in a record.

RecordControl's default name is "Data".

A RecordControl "mirrors" the contents of the controls it contains. Changing the values in the RecordControl will update the contents the controls. Similarly, when the contents of the controls are changed (by the user) the values in the RecordControl will be updated.

RecordControl is used by [Access1Control](<Access1Control.md>) and [AccessControl](<AccessControl.md>).

Most stdlib data entry controls are designed to work with RecordControl.

For example:

``` suneido
Window(Controller
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
    })
```

If you are writing your own control, here are the requirements for working with RecordControl:

-	`.Send(#Data)` in `New`
-	`.Send(#NoData)` in `Destroy`
-	`.Send(#NewValue, value)` when the user changes the value (commonly when the control loses the focus)
-	a `Get()` method that returns the value of the control
-	a `Set(value)` method that sets the control to the value (and sets dirty to false)
-	a `Dirty?(state = "")` method that returns whether the user has changed the value
-	a `Valid?()` method that returns whether the current value of the control is "valid"
-	optionally, a `SetValid(valid?)` method that sets whether the control is valid, (commonly used to color the control red)