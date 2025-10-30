### PopServer

``` suneido
(name = "POP Server", port = 110, exit = false)
```

This is a
[SocketServer](<SocketServer.md>)
that implements a simple POP server.

If exit is true, then closing the server window will exit from Suneido.
This is necessary if the PopServer is the only thing running in this instance of Suneido.

**Note:** PopServer is an abstract base class
- to use it you have to derive a concrete class
that defines:
`Authenticate(username, password)`
: Return true if the username and password are valid, false otherwise.  
There is a default definition that simply returns true.

`GetMessages(username)`
: Return an object (list) of the messages for the username
e.g. get them from a message store.

`Erase(i)`
: Delete the specified message.

`Complete()`
: This is called when the connection is closed.
For example, it can be used to complete a database transaction.
There is a default definition that does nothing.

For example:

``` suneido
MyPopServer

PopServer
    {
    GetMessages(username)
        {
        Database("ensure messages (when, from, to, msg) key(when)");
        t = Transaction(read:);
        q = t.Query("messages where to = " $ Display(username));
        .whens = Object();
        messages = Object();
        while (false isnt x = q.Next())
            {
            messages.Add(x.msg);
            .whens.Add(x.when);
            }
        t.Complete();
        return messages;
        }
    Erase(i)
        {
        t = Transaction(update:);
        t.Query("delete messages where when = " $ Display(.whens[i]));
        t.Complete();
        }
    }
```

See also: [SmtpServer](<SmtpServer.md>)