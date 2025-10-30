<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.Final

``` suneido
() => number
```

Returns the number of update transactions waiting to be finalized. If there are no outstanding transactions, this will be zero.

This is internal information primarily for debugging purposes.

Note: When client-server, this value comes from the server.

For example:

``` suneido
Database.Final()
    => 0
```