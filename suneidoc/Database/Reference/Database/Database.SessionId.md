<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.SessionId

``` suneido
(sessionid = "") => string
```

Returns the current session id. The session id is initially set to the IP address of the client. When standalone, the session id defaults to "127.0.0.1".

If an argument is supplied, the session id will be changed to it.

For example:

``` suneido
Database.SessionId()
    => "127.0.0.1"

Database.SessionId("fred")
    => "fred"
```

See also: [Database.Connections](<Database.Connections.md>)