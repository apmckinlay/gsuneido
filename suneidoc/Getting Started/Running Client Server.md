## Running Client Server

Okay, you've built this great application and now everyone in the office wants shared access to it.  No problem.  Up till now you've been running Suneido "locally" with the same instance of the executable acting as both database server and user interface client.  But Suneido can also run client-server with one instance of the executable acting as the database server, and multiple other instances, on other computers, acting as clients.

The first step is to create a text file called "server.go" containing:

``` suneido
Use("mylib");
```

**Note:** In client-server mode, it is the server that determines the set of libraries in use.

You can now start up Suneido in server mode with:

``` suneido
suneido -server server.go
```

**Note:** In practice you'd probably create a shortcut for starting the server.

To start up a client on the same machine run:

``` suneido
suneido -client 127.0.0.1 myset
```

where 127.0.0.1 is the special IP address for the current machine. This copy of Suneido will access all its data from the server via TCP/IP. i.e. this copy doesn't use a database file.

To start a client on another machine you will need to know the IP address of your server.  You can get this by running (on the server):

``` suneido
ipconfig
```

For example, if your server's IP address was 192.168.1.130 then to access it from another computer you would run:

``` suneido
suneido -client 192.168.1.130 myset
```

Again, in practice you would probably create a shortcut for your users to run this command.  And to avoid having to copy suneido.exe to all your client machines, you could put it in a shared location like \\server\shared and then run:

``` suneido
\\server\shared\suneido -client 192.168.1.130 myset
```

This also means that you only have to update one copy of the executable when there is a new version.

By default Suneido uses TCP/IP port 3147.  If necessary you can specify a
different port on both the server and clients:

``` suneido
suneido -port 1234 -server server.go
```

``` suneido
suneido -port 1234 -client 192.168.1.130 myset
```

You'll also need to use a different port if you want to run more than one Suneido server on the same computer.  (Each server on a given computer needs its own port.)

The one drawback with this is that there will only be one persistent state. For example, if one user maximizes his window, then everyone will get maximized windows. This usually isn't what you want. In order to have a different persistent state for each user, it's necessary for users to *login* to identify themselves. Here's a simple Login:

``` suneido
function (set)
    {
    if (false is user = Ask('Name', 'Login'))
        Exit()
    else
        {
        Suneido.User = Suneido.User_Loaded = user
        PersistentWindow.Load(set)
        }
    }
```

Then you can run (from the command line or a shortcut:

``` suneido
suneido Login("myset")
```

or as a client:

``` suneido
suneido -client 192.168.1.130 Login("myset")
```

**Note:** This Login function is very simple. Other applications might require a more sophisticated login with, for example, user validation, passwords, and [Database.Auth](<../Database/Reference/Database/Database.Auth.md>).

**Note:** The server will automatically disconnect clients that are inactive (make no requests) for two hours.

In production you would probably want to run your server as a Windows service. See [Running as a Service](<../Introduction/Running as a Service.md>)