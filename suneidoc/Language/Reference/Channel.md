<div style="float:right"><span class="builtin">Builtin</span></div>

### Channel

``` suneido
(size = 4) => instance
```

A channel is a Go channel of Suneido values.

Methods:
`Send(value)`
: Sends a value to the channel. If the channel size is 0 (unbuffered) this will block until another thread calls Recv. Otherwise it will only block if the channel buffer is full. Send to a closed channel will throw an exception.

`Recv() => value`
: Gets a value from the channel. Will block if no value is available. Recv from a closed channel will return the channel itself.

`Recv2(channel2) => object(0|1 [, value])`
: Get value from whichever of the two channels has one available first. If both channels have values available immediately, it will choose randomly. The returned object has one or two items. The first is which channel returned the result, 0 or 1. The second is the value received, unless the channel is closed. Recv2 is useful to monitor a "done" channel in addition to a data channel.

`Close()`
: Closes the channel. Normally the sender should close the channel. Do not Close in one thread at the same time as another thread is Send'ing. It is not required to Close channels unless you need to unblock receivers. Closing a closed channel will throw an exception.

Send, Recv, and Recv2 will block for a maximum of 10 seconds and then throw a timeout.

Some possible configurations:

-	One thread sending to a channel and one thread receiving  i.e. a pipeline
-	One thread sending to a channel and multiple threads receiving i.e. fan out
-	Multiple threads sending to a channel and one thread receiving i.e. fan in


Channel is a thin wrapper around a Go channel with the addition of blocking timeouts. For more information see the Go documentation. For example: [Go Concurrency Patterns: Pipelines and cancellation](<https://go.dev/blog/pipelines>)

**WARNING**: Concurrency is hard. Only use it if the benefits are substantial. Never take it lightly.


See also:
[Mutex](<Mutex.md>),
[object.CompareAndSet](<Object/object.CompareAndSet.md>),
[Synchronized](<Synchronized.md>),
[class.Synchronized](<Class/class.Synchronized.md>),
[Thread](<Thread.md>),
[WaitGroup](<WaitGroup.md>)
