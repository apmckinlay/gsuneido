### [Controller](<../Controller.md>) - Methods
`On_About_Suneido()`
: Displays the About Suneido window.   
Allows "About Suneido" to be placed on the menu without having to define it.

`On_Users_Manual()`
: Opens the Suneido help (suneidoc book) in a persistent window.   
Allows "Users Manual" to be placed on the menu without having to define it.

`Redir(message, target = "focus")`
: Redirects the specified message to the specified control or function, or to the control that has the focus.  This is a shortcut to receiving the message and then re-sending it.   
By default, the standard editing commands are redirected to the focus control.  (i.e. Cut, Copy, Paste, Undo, Redo)