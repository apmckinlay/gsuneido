## Message Loop

Once Suneido has run Init (see [Startup](<Startup.md>)) 
it enters a message loop that processes Windows messages. 
This is a fairly standard message loop with 
GetMessage, TranslateAccelerators, IsDialogMessage, TranslateMessage, DispatchMessage.

Dialogs can run their own message loop with [MessageLoop](<../User Interfaces/Reference/MessageLoop.md>)