### SmtpServer

``` suneido
(name = "SMTP Server", port = 25, exit = false)
```

This is a [SocketServer](<SocketServer.md>) that implements a simple SMTP server.

If exit is true, then closing the server window will exit from Suneido. This is necessary if the SmtpServer is the only thing running in this instance of Suneido.

**Note:** SmtpServer is an abstract base class - to use it you must derive a concrete class that defines:
`Recipient?(rcpt)`
: Return true if rcpt is valid, false otherwise.  
There is a default definition that simply returns true.

`Send(from, to, msg)`
: Process messages e.g. store them in a message table or forward them somewhere.

For example:

``` suneido
MySmtpServer

SmtpServer
    {
    Send(from, to, msg)
        {
        Database("ensure messages (when, from, to, msg) key(when)");
        t = Transaction(update:);
        q = t.Query("messages");
        q.Output(Record(when: Timestamp(), from: from, to: to, msg: msg));
        t.Complete();
        }
    }
```

See also: [Pop3Server](<Pop3Server.md>)