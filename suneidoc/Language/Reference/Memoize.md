### Memoize

Abstract base class for caching the results of a function.

Derived classes must define Func and may optionally define Init and CacheSize (default 100).

For example, a recursive fibonacci function can be improved by caching results. (In practice an iterative version would be better.)

``` suneido
Fibonacci

function(n)
    {
    return (n <= 1) ? n : (Fibonacci(n - 1) + Fibonacci(n - 2))
    }
```

On my machine this version takes 3 seconds to calculate Fibonacci(32) whereas the following version takes about 1.5 ms.

``` suneido
Memoize
    {
    CacheSize: 4
    Func(n)
        {
        return (n <= 1) ? n : Fibonacci(n - 1) + Fibonacci(n - 2)
        }
    }
```

Memoize creates an [LruCache](<LruCache.md>) to cache the result and stores it in the global Suneido object.

If Func() argument(s) size are too big, you can specify HasArgs? member to be true on the derived classes, so [LruCache](<LruCache.md>) will hash argument to save memory

Memoize also defines a ResetCache method that will reset the associated LruCache.


See also:
[Compose](<Compose.md>),
[Curry](<Curry.md>),
[MemoizeSingle](<MemoizeSingle.md>),
[object.Any?](<Object/object.Any?.md>),
[object.Drop](<Object/object.Drop.md>),
[object.Every?](<Object/object.Every?.md>),
[object.Filter](<Object/object.Filter.md>),
[object.FlatMap](<Object/object.FlatMap.md>),
[object.Fold](<Object/object.Fold.md>),
[object.Map](<Object/object.Map.md>),
[object.Map!](<Object/object.Map!.md>),
[object.Map2](<Object/object.Map2.md>),
[object.Nth](<Object/object.Nth.md>),
object.PrevSeq,
[object.Reduce](<Object/object.Reduce.md>),
[object.Take](<Object/object.Take.md>),
[object.Zip](<Object/object.Zip.md>)
