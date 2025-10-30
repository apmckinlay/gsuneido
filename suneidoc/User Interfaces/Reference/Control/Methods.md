### [Control](<../Control.md>) - Methods
`AlertError(title, message)`
: Does:  
`Alert(message, title, .Window.Hwnd, MB.ICONERROR)`

`AlertInfo(title, message)`
: Does:  
`Alert(message, title, .Window.Hwnd, MB.ICONINFORMATION)`

`AlertWarn(title, message)`
: Does:  
`Alert(message, title, .Window.Hwnd, MB.ICONWARNING)`

`Construct(object or @object)`
: Creates an instance of a control, normally accessed as .Construct.   
If the name of the control starts with a lower case letter, it is assumed to be a field name and the control is looked up in the data dictionary and merged with the supplied arguments.  Otherwise, "Control" is appended onto the name.  Then Construct is used to create the instance.   
Automatically copies the following members to the new instance as public members (capitalized) in Control.New (i.e. prior to the control's New):  
`xmin, ymin, top, xstretch, ystretch, name`  
If the constructed control has a name, then the control is assigned to a member of its parent. This allows the parent to access its controls by name.

`FindControl(name)`
: Returns a reference to the named control or false if not found. Uses .Children() to search the "tree" of controls. Custom containers must implement Children if you want FindControl to be able to find their child controls. Most controls have a default name: (e.g. ButtonControl's default name is "Button") but you will need to override these with unique names to use FindControl. Note: When you include a field by name (i.e. to use the control from its Field_ definition) then the field name will be automatically added to the control.

`Send(message [, argument ...])`
: Send a message to a control's controller. Automatically  adds a *source* argument to the message.