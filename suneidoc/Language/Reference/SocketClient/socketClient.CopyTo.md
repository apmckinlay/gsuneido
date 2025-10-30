<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketClient.CopyTo

``` suneido
(dest, nbytes = false) => number
```

Reads from the source and writes to the destination, until nbytes or eof. The source and destination can be open files, pipes, or sockets. **CopyTo** is more efficient than using Read and Write because it reuses the same buffer and so it does much less memory allocation.

Returns the number of bytes copied.

If nbytes is specified and it hits eof before reading the full amount then it is [return throw](<../../Statements/return.md>). i.e. if the result is not used, it will throw an exception.

**Warning**: If reading from a pipe or socket, and not specifying nbytes, CopyTo will not return until the source pipe or socket is closed. A socket might time out, but a pipe will not.


See also:
[file.CopyTo](<../File/file.CopyTo.md>),
pipe.CopyTo,
[socketServer.CopyTo](<../SocketServer/socketServer.CopyTo.md>)
