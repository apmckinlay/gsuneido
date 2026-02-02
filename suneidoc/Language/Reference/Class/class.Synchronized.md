<div style="float:right"><span class="builtin">Builtin</span></div>

#### class.Synchronized

``` suneido
(block) => value
```

.Synchronized locks the current class's internal mutex, runs the block, and then unlocks the mutex. It returns whatever the block returns. It will timeout and throw an exception if it is unable to obtain the lock within 10 seconds.

For more information see the documentation for [Mutex](<../Mutex.md>). .Synchronized is equivalent to Mutex except that you do not need to create or store the mutex yourself.

**Note**: .Synchronized is <u>not</u> reentrant. If code calls .Synchronized while in .Synchronized (in the same class) it will deadlock (and then timeout).

**Note**: .Synchronized cannot be used in nested classes. (Because of the way it locates the class for the currently executing function.)

**WARNING**: Concurrency is hard. Only use it if the benefits are substantial. Never take it lightly.


See also:
[Channel](<../Channel.md>),
[Mutex](<../Mutex.md>),
[object.CompareAndSet](<../Object/object.CompareAndSet.md>),
[Pipe](<../Pipe.md>),
[Synchronized](<../Synchronized.md>),
[Thread](<../Thread.md>),
[WaitGroup](<../WaitGroup.md>)
