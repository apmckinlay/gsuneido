<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.Cursors

``` suneido
() => number
```

Returns the number of open cursors for the current thread.

**Note:** This only works client-server. It will always return 0 when running standalone.

For example:

``` suneido
Database.Cursors()
    => 0
```