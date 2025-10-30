### Writing a Controller

A [Controller](<../Reference/Controller.md>) wraps a control (or group of controls) in order to manage them and receive messages they Send.

**`MyEditorControl`**
``` suneido
Controller
    {
    Title: "Simple Editor"
    New()
        {
        e = .Vert.Editor
        .Redir('On_Undo', e)
        .Redir('On_Redo', e)
        .Redir('On_Cut', e)
        .Redir('On_Copy', e)
        .Redir('On_Paste', e)
        .status = .Vert.Status
        }
    Controls:
        (Vert
            (Toolbar Cut Copy Paste "" Undo Redo)
            Scintilla
            Statusbar)
    Commands:
        (
        (Undo,              "Ctrl+Z",   "Undo the last action")
        (Redo,              "Ctrl+Y",   "Redo the last action")
        (Cut,               "Ctrl+X",   "Cut the selected text to the clipboard")
        (Copy,              "Ctrl+C",   "Copy the selected text to the clipboard")
        (Paste,             "Ctrl+V",   "Insert the contents of the clipboard")
        )
    Menu:
        (
        ("&File",
            "&Close")
        ("&Edit",
            "&Undo", "&Redo", "",
            "Cu&t", "&Copy", "&Paste")
        )
    Status(status)
        {
        .status.Set(status)
        }
    }
```

Note:

-	inherit from Controller
-	specify a Title *
-	optionally redirect commands (menu or toolbar) to controls
-	specify the control(s)
-	define commands (menu or toolbar)
-	define a menu if you want one *
-	handle messages that controls Send e.g. Status from Scintilla


* Title and Menu will only be used if this is the top level control in a Window or Dialog

Note: Any messages that the Controller does not handle are dropped. If you want unhandled messages to propagate to another Controller that contains this one, use [PassthruController](<../Reference/PassthruController.md>) instead.