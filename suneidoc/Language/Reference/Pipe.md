<div style="float:right"><span class="builtin">Builtin</span></div>

### Pipe

``` suneido
() => reader, writer
```

Creates a synchronous in-memory pipe (wrapping a Go io.Pipe). Returns two values: a reader and a writer.

Writes to the writer are matched with reads from the reader. Writes and reads both block until the other side is ready. This makes pipes useful for coordinating between threads.

**Note:** This requires two threads. Trying to read and write from the same thread will deadlock.

**Note:** This is different from the operating system pipes used by [RunPiped](<RunPiped.md>).

#### Reader Methods:

`Read(n) => string | false`
: Reads up to n bytes from the pipe. Returns a string of the bytes read, or false if the pipe is closed (EOF). Maximum read size is 64kb. Blocks until data is available or the pipe is closed.

`CopyTo(dest, nbytes = false) => ncopied`
:	Up to nbytes (or everything) is read from this and written to a suitable destination. 
	e.g. [File](<File.md>),
	[Pipe](<Pipe.md>),
	[HttpClient2](<HttpClient2.md>),
	[HttpServer](<HttpServer.md>),
	[RunPiped](<RunPiped.md>),
	[SocketClient](<SocketClient.md>),
	[SocketServer](<SocketServer.md>).
	Prefer `CopyTo` when applicable since it is more efficient.

`Close()`
: Closes the reader side of the pipe. No further reads will be possible.

#### Writer Methods:

`Write(s)`
: Writes a string to the pipe. Blocks until the reader is ready to receive.

`Close()`
: Closes the writer side of the pipe. This signals EOF to the reader. Normally the writer should close the pipe when done.

Pipe is a thin wrapper around Go's io.Pipe. For more information see the Go documentation for [io.Pipe](<https://pkg.go.dev/io#Pipe>).

**WARNING**: Concurrency is hard. Only use it if the benefits are substantial. Never take it lightly.


See also:
[Channel](<Channel.md>),
[Mutex](<Mutex.md>),
[object.CompareAndSet](<Object/object.CompareAndSet.md>),
[Synchronized](<Synchronized.md>),
[class.Synchronized](<Class/class.Synchronized.md>),
[Thread](<Thread.md>),
[WaitGroup](<WaitGroup.md>)
