<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketServer.SetTimeout

``` suneido
(seconds)
```

Sets the timeout for the socket.

Note: Writes are buffered by the operating system so they will not usually timeout unless they are larger than the buffer.