<div style="float:right"><span class="builtin">Builtin</span></div>

### WaitGroup

``` suneido
() => instance
```

A WaitGroup is used to wait for a group of threads to finish.

Methods:
`Thread(callable, name = "")`
: Equivalent to 
[Thread](<Thread.md>) but automatically does Add() and ensures Done() is called when the thread ends.

`Add(inc = 1)`
: Adds one to the wait group counter. (Should be called before starting the thread.)

`Done()`
: Decrements the wait group counter. It will throw an exception if the counter becomes negative.

`Wait(secs = 10) => true or "WaitGroup: timeout"`
: Blocks for up to timeoutSeconds until the wait group counter becomes zero. It is return-throw, so if the result is not used and it times out, it will throw an exception.

**Note**: Using Thread is preferable to manually doing Add() and Done().

For example:

``` suneido
wg = WaitGroup()
wg.Thread(Task1)
wg.Thread(Task2)
wg.Thread(Task3)
wg.Wait()
```

**Warning**: If you call Wait from the main UI thread, it will block any Print's in threads.


See also:
[Channel](<Channel.md>),
[Mutex](<Mutex.md>),
[object.CompareAndSet](<Object/object.CompareAndSet.md>),
[Pipe](<Pipe.md>),
[Synchronized](<Synchronized.md>),
[class.Synchronized](<Class/class.Synchronized.md>),
[Thread](<Thread.md>)
