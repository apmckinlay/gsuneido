<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/PersistentWindow/Methods">Methods</a></span></div>

### PersistentWindow

``` suneido
(control = false, master? = false, stateobject = #(), 
    newset = false, exitOnClose = false)
```

PersistentWindow is used to create a window that needs to save and restore 
its state each time it is run.  The control is responsible for setting and 
retrieving its state via SetState and GetState methods. The top level window 
saves its size and position. PersistentWindow can also be used to create an 
application's window which can be run separately from the IDE.

The first argument, **control**, is the control specification.  This is 
specified the same way [Window](<Window.md>)
's control argument is specified.

If **master?** is true, then the resulting window will be the master 
window for the current instance of suneido.exe, which means closing this window 
will close suneido.exe (and all other windows belonging to that instance of suneido.exe).

**stateobject** contains the necessary information to restore the window
to its previous state.  It is not necessary to include this argument when creating
a persistent window the first time, it gets used by PersistentWindow when it loads
a previously defined persistent window set (it creates instances of itself when 
loading a persistent window set).

**newset** is used to create a new persistent window set.  This is 
useful for setting up an application window.  This argument should be a 
string representing the name of the new set.  Windows created with this 
option are always master windows, see master? argument above.

For example, to create a persistent set that just loads the Suneido 
help, run the following from the workspace: 

``` suneido
PersistentWindow(#(Book suneidoc), newset: "help")
```

This will close all of your current IDE windows and create the persistent
window.  Closing this window will close suneido.exe.  You can then run this help
window from the command line:

``` suneido
suneido help
```

This same approach can be used to set up shortcuts for running a Suneido application.

See also: 
[Window](<Window.md>)