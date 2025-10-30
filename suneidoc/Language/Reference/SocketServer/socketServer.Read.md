<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketServer.Read

``` suneido
(nbytes = false) => string or false
```

Returns the next **nbytes**, or the rest of the data if **nbytes** is false, or false if no data is available and the connection is closed.

If the connection is closed and partial data has been received, it will be returned. i.e. the returned string may not always be **nbytes** long. The next call to Read will return false.

Throws an exception if the timeout is exceeded.

**Note:** Prior to BuiltDate 20241219, nbytes had to be specified.