<div style="float:right"><span class="builtin">Builtin</span></div>

#### Thread.Name

``` suneido
(string = false) => name
```

With no argument, [Thread.Name](<Thread.Name.md>)() returns the current thread name.

With a string argument, Thread.Name(string) sets the "extra" information on the current thread's name.

Threads have a built-in thread name which is not altered.

-	Threads started with 
	[Thread](<../Thread.md>) will be named "Thread-#"


The extra information specified by Thread.Name(string) will be appended to the built-in name. (Replacing any existing extra information.)

[Thread.List](<Thread.List.md>)() will show the full thread names, including any extra information set with Thread.Name(string).

The thread name can also be passed to the initial [Thread](<Thread.md>) call.