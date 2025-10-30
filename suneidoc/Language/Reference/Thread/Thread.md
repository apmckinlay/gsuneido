<div style="float:right"><span class="builtin">Builtin</span></div>

### Thread

``` suneido
(callable, name = false)
```

Thread creates a new thread to run the specified function or class. Thread returns immediately, without waiting for the new thread.

If a name is supplied it has the same effect as doing [Thread.Name](<Thread.Name.md>).

For example:

``` suneido
f = function ()
    {
    for (i = 0; i < 10; ++i)
        {
        Beep()
        for (j = 0; j < 1000000; ++j)
            {} // delay
        }
    }
Thread(f)
```

This will beep 10 times, while allowing you to continue to work.

[SocketServer](<../SocketServer.md>) creates threads for each accepted connection.

**Note**: Code run by Thread should not access the Windows GUI. i.e. it should not get or set anything to do with windows or controls. A Thread can use [Defer](<../../../User Interfaces/Reference/Defer.md>) to run code on the main thread to update the GUI.