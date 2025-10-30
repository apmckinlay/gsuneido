## Responding to a Button

**Category:** User Interface

**Problem**

You want to do something when the user clicks on a button.

**Ingredients**

[ButtonControl](<../User Interfaces/Reference/ButtonControl.md>), 
[Controller](<../User Interfaces/Reference/Controller.md>), 
[FieldControl](<../User Interfaces/Reference/FieldControl.md>), 
[StaticControl](<../User Interfaces/Reference/StaticControl.md>), 
[Alert](<../User Interfaces/Reference/Alert.md>)

**Recipe**

Enter this in a library (e.g. mylib) as My_Button:

``` suneido
Controller
    {
    Controls:
        (Vert
            (Static 'Enter your name')
            Field
            (Button Hello)
            )
    On_Hello()
        {
        Alert("Hello " $ .Vert.Field.Get())
        }
    }
```

You can then run it directly from Library View or from the WorkSpace with:

``` suneido
My_Button()
```

or:

``` suneido
Window(My_Button)
```

Controls like Button use .Send to send "messages" to their Controller. To handle/respond to these messages you have to create a Controller and within it, define a method corresponding to the
message that the control is Send'ing. In this case, ButtonControl sends a message of "On_" $ button-name and our button is called "Hello" so our method is "On_Hello". In this example, the On_Hello
method simply displays an Alert - but it could, of course, do other things.

**Discussion**

The example code also demonstrates one way to access the contents of a field control - a better alternative if you have a lot of fields is to use a RecordControl.

If you only want your controller to catch some messages and to pass the rest on to any containing controller, use PassthruController instead of Controller.

**See Also**

[Using RecordControl](<Using RecordControl.md>), 
[Determining the Source of Messages](<Determining the Source of Messages.md>)