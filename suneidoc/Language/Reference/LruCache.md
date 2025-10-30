<div style="float:right"><span class="builtin">Builtin</span></div>

### LruCache

``` suneido
(getfn, cache_size = 10, hashArgs? = false, okForResetAll? = true)
```

#### Methods
`Get(key) => value`
: If key is present in the cache, it is returned. Otherwise getfn is called with key as an argument to calculate or fetch it and the result is saved in the cache, evicting the least recently used entry if the cache is full.

`GetN(@args) => value`
: Similar to Get, but uses all the arguments as the key, and passes the same arguments to getfn

`GetN1(args) => value`
: Similar to GetN, but takes all arguments as one object, and passes the same arguments to getfn

`GetMissRate() => number`
: Returns the number of "misses" (when getfn had to be called) divided by the number of gets. Lower is better.

`lrucache.Reset()`
: Clear the cache.

`LruCache.ResetAll()`
: Deletes all Suneido members that are LruCache's

See also: [Memoize](<Memoize.md>)