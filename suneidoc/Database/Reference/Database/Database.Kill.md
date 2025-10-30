<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.Kill

``` suneido
(session_id) => number
```

Kills any sessions with the specified session_id. Returns the number of sessions killed.

For example:

``` suneido
Database.Kill("192.168.1.143")
    => 1
```

This is useful if Database.Connections shows a connection but that machine is not running a client. (This can happen if the client is not closed properly.)

This only works client-server. When running standalone it will always return 0.

See also:
[Database.Connections](<Database.Connections.md>),
[Database.SessionId](<Database.SessionId.md>)