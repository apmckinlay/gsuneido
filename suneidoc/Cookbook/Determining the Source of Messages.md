## Determining the Source of Messages

**Category:** User Interface

**Problem**

You have several controls that send the same message. When your controller receives the message, how do you determine which of the controls sent it?

**Ingredients**

[Controller](<../User Interfaces/Reference/Controller.md>), 
[FieldControl](<../User Interfaces/Reference/FieldControl.md>), 
[Alert](<../User Interfaces/Reference/Alert.md>)

**Recipe**

Enter this in a library (e.g. mylib) as My_Controller:

``` suneido
Controller
    {
    Controls: (Vert
        (Field name: one)
        (Field name: two))
    New()
        {
        .one = .Vert.one
        }
    NewValue(value, source)
        {
        if (source is .one)
            {
            Alert('one changed')
            .one.Dirty?(false)
            }
        }
    }
```

You can then run it directly from Library View or from the WorkSpace with:

``` suneido
My_Controller()
```

or:

``` suneido
Window(My_Controller)
```

**Discussion**

In this example, we define a controller with two Field controls. Field controls send a **NewValue** message when they lose the focus and are *dirty* (i.e. the user has changed their
contents).

Messages include an optional final "source" argument which is a reference to the control that sent the message. We can use this source argument to determine which of the controls sent the message.
In this example, if the source is the field named "one" then we give an Alert.

**Note:** The receiver of a NewValue message is responsible for clearing the dirty state of the control.

Although we could just compare the source to .Vert.one directly, instead we follow the standard practice of obtaining a reference to the control in the New. The advantage of this is that if you
change the view layout, you only have to change the *path* to the control in one place (in New) instead of throughout your code.

Notice that in order to refer to the fields individually, we had to give them *names*. Controls have default names, but when you have multiple controls of the same type it is usually
necessary to specify different names.

**Note:** Because the source argument is passed as a named argument (so that it is optional on the receiver) it must be called "source". You cannot use a different name.

**See Also**

[Responding to a Button](<Responding to a Button.md>), 
[Using RecordControl](<Using RecordControl.md>)