### Construction

Controls are created using Control.Construct.  Control.New sets the following members:

.Window
: The top level 
[Window](<../Reference/Window.md>) or 
[Dialog](<../Reference/Dialog.md>).

.Parent
: The containing Suneido control.

.HwndCtrl
: The containing Win32 window.

Note: .HwndCtrl and .Parent are not necessarily the same control.

After creating the control, Construct will automatically set some of it's members from the control specification:

xmin => Xmin

ymin => Ymin

xstretch => Xstretch

ystretch => Ystretch

name => Name

This is a shortcut so controls don't all have to handle these arguments.  However, since this isn't done till after the control is created, the members will not be available during New.  If these members are required by New, then it must handle them as arguments.

After all the controls are contructed, **Startup** is called on each control, top down. This is useful to perform setup that can not be done in New because construction is not complete at that point.