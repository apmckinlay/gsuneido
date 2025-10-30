<div style="float:right"><span class="builtin">Builtin</span></div>

### Mutex

``` suneido
() => instance

instance.Do(block) => result
```

A Mutex instance has a single method, Do, which runs a function or block exclusively. For a given mutex instance, if multiple threads call Do at the same time, only one at a time will run their block. The others will wait (block). mutex.Do returns the result of the block. It will timeout and throw an exception if it is unable to obtain the lock within 10 seconds.

**Note**: .Synchronized is <u>not</u> reentrant. If code calls .Synchronized while in .Synchronized (in the same class) it will deadlock (and then timeout).

Most built-in operations in Suneido are "atomic". However, if you need to do several operations atomically, then you can use a Mutex. For example:

``` suneido
if Suneido.Member?(#foo)
	++Suneido.foo
```

This is NOT safe for concurrent use because another thread could remove Suneido.foo in between Member? and ++. You can use a mutex to ensure this cannot happen.

``` suneido
mutex.Do()
	{
	if Suneido.Member?(#foo)
		++Suneido.foo
	}
```

However, this means ALL the code that modifies Suneido.foo must use the (same) mutex, even if that code is a single operation that by itself is atomic.

``` suneido
mutex.Do()
	{
	Suneido.Delete(#foo)
	}
```



Single "read" operations do not require using the mutex.

``` suneido
Print(Suneido.foo)
```

However, if the mutex is guarding several variables, and you read more than one, then you need to use the mutex to ensure consistent values between the two.

``` suneido
mutex.Do()
	{
	Print(Suneido.foo)
	Print(Suneido.bar)
	}
```

**Note**: Mutex is not reentrant i.e. cannot be nested. If you call mutex.Do from within mutex.Do it will block and eventually time out (after 10 seconds).

**Note**: You only need to use Mutex when there are multiple threads. Single threaded code does not need to use Mutex.

The code run by mutex.Do should be small and fast. It should not call anything that could be slow e.g. user interface, network, or file system.

A mutex should guard particular data. Don't share a mutex for unrelated uses.

Mutex replaces Synchronized (which was basically a single global shared mutex).

**WARNING**: Avoid nesting different mutexes. Unless you guarantee they are always nested in the same order, you risk deadlock (which means they will time out). For example:

``` suneido
thread1:
	mutex1.Do()
		{
		mutex2.Do() // POTENTIAL DEADLOCK
			{ ... }
		}

thread2:
	mutex2.Do()
		{
		mutex1.Do() // POTENTIAL DEADLOCK
			{ ... }
		}
```

A Mutex is not a Go mutex. It is implement with a Go channel to support timeouts.

**WARNING**: Concurrency is hard. Only use it if the benefits are substantial. Never take it lightly.


See also:
[Channel](<Channel.md>),
[object.CompareAndSet](<Object/object.CompareAndSet.md>),
[Synchronized](<Synchronized.md>),
[class.Synchronized](<Class/class.Synchronized.md>),
[Thread](<Thread.md>),
[WaitGroup](<WaitGroup.md>)
