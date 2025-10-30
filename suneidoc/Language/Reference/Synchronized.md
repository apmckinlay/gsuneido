<div style="float:right"><span class="builtin">Builtin</span><span class="deprecated">Deprecated</span></div>

### Synchronized

``` suneido
(callable)
```

Synchronized is **deprecated**. Please use [class.Synchronized](<Class/class.Synchronized.md>) or [Mutex](<Mutex.md>) instead.

Runs the specified callable (e.g. block, function, method) without yielding to any other threads. This is useful to prevent concurrent threads from interfering with each other.

Synchronized uses a single global reentrant lock.

**Note**: It is important to keep Synchronized blocks as small and fast as possible so they do not cause delays for other threads.

For example: (using a block)

``` suneido
Synchronized()
    {
    --Suneido.incoming
    ++Suneido.outgoing
    }
```

Making this synchronized prevents other threads from seeing the Suneido members in an inconsistent state.

See also:
[Thread](<Thread.md>)