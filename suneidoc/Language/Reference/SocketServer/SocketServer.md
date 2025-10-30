<div style="float:right"><span class="builtin">Builtin</span></div>

#### SocketServer

``` suneido
(name = .Name, port = .Port)
```

SocketServer is the base class for TCP/IP socket stream servers.

Classes derived from Server must define a Run method.

Calling the server class starts it listening for connections. The **name** and **port** may be passed in or supplied by the class. An instance of the class is created to service each connection.

Each connection instance runs in a separate thread.

Note: Returning from Run closes that connection (but not the server).

Here is a simple example:

``` suneido
TestServer

SocketServer
    {
    Port: 1234
    Run()
        {
        .Writeline("hello")
        while (false isnt (req = .Readline()))
            if (req is "quit")
                break
            else
                .Writeline("don't know how to " $ req)
        }
    }
```

Which can then be run with: `TestServer()`

You can pass additional arguments to your server by defining New. New will be called once when you start the server (not for each connection). For example:

``` suneido
    New(.greeting)
        { 
        }
    Run()
        {
        .Writeline(.greeting)
    ...
```

The way this works is that when you start the socket server, it creates a "master" instance, calling New. Then for each connection it makes a copy of this master instance.

When it passes the arguments to New, it ensures that the SocketServer's arguments (name and port) are named. For example:

``` suneido
TestServer("myname", 9000, "extra argument", another: 123)
```

will become:

``` suneido
New("extra argument", another: 123, name: "myname", port: 9000)
```

So the SocketServer's arguments (name, port) are available if you want them. But since they are named you do not have to receive them. For example, your New could be either of:

``` suneido
New(extra, another)
```
or:
``` suneido
New(extra, another, name, port)
```

If you do receive them, they have to be the same names (name, port).

If you are not passing name and port then you will need to name your own arguments. e.g. `MySocketServer(myarg: ...)` since if you just did `MySocketServer(myarg)` it would be taken as name, and your New would not receive its argument.

If the SocketServer class contains a method named "Killer" it will be called at startup with a "killer" object that can be used to stop the SocketServer. The method will need to save this killer object somewhere e.g. in the Suneido object. **Note:** Stopping the server does not stop any threads for existing connections. If required, that has to be handled by the application code e.g. setting a flag that the connection code checks.