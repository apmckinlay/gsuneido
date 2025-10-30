<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.TempDest

``` suneido
() => number
```

Returns the current amount of temporary index space in use. If all queries are closed properly, this should be zero.

Note: When client-server, this value comes from the server.

For example:

``` suneido
Database.TempDest()
    => 0
```