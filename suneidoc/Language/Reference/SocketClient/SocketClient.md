<div style="float:right"><span class="builtin">Builtin</span></div>

#### SocketClient

``` suneido
(address, port, timeout = 60, timeoutConnect = 0) => socketClient
(address, port, timeout = 60, timeoutConnect = 0) {|socketClient| ... }
```

Opens a socket connection to the specified address and port. The address can be either a name (e.g. www.suneido.com) or dotted numeric format (e.g. 127.0.0.1).

**timeout** is the maximum number of seconds (default 60 seconds) that the connection will wait for read's or write's. If the timeout is exceeded, an exception will be thrown.

**timeoutConnect** is the maximum number of seconds to wait to make the initial connection. A value of 0 (the default) uses the system settings (usually quite long e.g. 60 seconds) Decimals are allowed (e.g. timeoutConnect: .5)

For example:

``` suneido
sc = SocketClient("192.168.1.130", 110)
Print(sc.Readline())
sc.Close()
```

It is  safer to use the block form of SocketClient to ensure the socket gets closed even if there are exceptions.

``` suneido
SocketClient("192.168.1.130", 110)
    { |sc|
    Print(sc.Readline())
    }
```

See also: 
[SocketServer](<../SocketServer.md>),
[Thread](<../Thread.md>)