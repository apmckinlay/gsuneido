<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketServer.CopyTo

``` suneido
(dest, nbytes = false) => number
```
Up to nbytes (or everything) is read from this and written to a suitable destination. 
e.g. [File](<../File.md>),
[Pipe](<../Pipe.md>),
[HttpClient2](<../HttpClient2.md>),
[HttpServer](<../HttpServer.md>),
[RunPiped](<../RunPiped.md>),
[SocketClient](<../SocketClient.md>),
[SocketServer](<../SocketServer.md>).
`CopyTo` is more efficient than using Read and Write because it reuses the same buffer 
and so it does much less memory allocation.

Returns the number of bytes copied.

If nbytes is specified and it hits eof before reading the full amount then it is [return throw](<../../Statements/return.md>). i.e. if the result is not used, it will throw an exception.

**Warning**: If reading from RunPiped or a socket, and not specifying nbytes, CopyTo will not return until the source pipe or socket is closed. A socket might time out, but a pipe will not.


See also:
[file.CopyTo](<../File/file.CopyTo.md>),
[HttpClient2](<../HttpClient2.md>),
HttpsClient,
[HttpServer](<../HttpServer.md>),
pipe.CopyTo,
[runPiped.CopyTo](<../RunPiped/runPiped.CopyTo.md>),
[socketClient.CopyTo](<../SocketClient/socketClient.CopyTo.md>)
