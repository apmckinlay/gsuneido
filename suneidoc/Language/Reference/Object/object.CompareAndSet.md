<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.CompareAndSet

``` suneido
(member, newval [, oldval]) => true or false
```

CompareAndSet is an atomic equivalent of:

``` suneido
if ob[member] is oldval
	{
	ob[member] = newval
	return true
	}
else
	return false
```

If oldval is not specified then the member must not exist.

A common use is to ensure something is done a single time across multiple threads. Multiple threads can run this and CompareAndSet will only succeed on one of them. (Assuming Suneido.done is initialized to false.

``` suneido
if Suneido.CompareAndSet(#done, false, true)
	something() // will only be done once
```

NOTE: CompareAndSet is only needed when there are multiple threads.


See also:
[Channel](<../Channel.md>),
[Mutex](<../Mutex.md>),
[Pipe](<../Pipe.md>),
[Synchronized](<../Synchronized.md>),
[class.Synchronized](<../Class/class.Synchronized.md>),
[Thread](<../Thread.md>),
[WaitGroup](<../WaitGroup.md>)
