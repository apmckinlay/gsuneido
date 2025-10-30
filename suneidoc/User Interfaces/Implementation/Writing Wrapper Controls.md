### Writing Wrapper Controls

A "wrapper" control is a Suneido control that "wraps" a Windows control. Some examples of wrapper controls are: [ButtonControl](<../Reference/ButtonControl.md>), [StaticControl](<../Reference/StaticControl.md>), and [EditControl](<../Reference/EditControl.md>).

Wrapper controls normally derive from either Hwnd or WndProc if they need to respond to WM_ window messages.

#### Guidelines

-	Write methods for each SendMessage. Name these messages based on the Windows API message constant. e.g. the method for EM_REPLACESEL would be called ReplaceSel.
-	Keep the wrapper "thin". If you want to add behavior consider doing it in a derived class. For example, EditControl is the basic wrapper with FieldControl and EditorControl adding additional functionality


See also: [Writing a Hwnd Control](<Writing a Hwnd Control.md>), [Writing a WndProc Control](<Writing a WndProc Control.md>).