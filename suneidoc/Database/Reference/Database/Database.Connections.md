<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.Connections

``` suneido
() => object
```

Returns a list of the current connections (session ids) for the database server.

Note: This list will be empty when running standalone.

For example:

``` suneido
Database.Connections()
    => #("127.0.0.1")
```

See also: [Database.SessionId](<Database.SessionId.md>)