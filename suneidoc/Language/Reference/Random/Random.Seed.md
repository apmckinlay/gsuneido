<div style="float:right"><span class="builtin">Builtin</span></div>

### Random.Seed

``` suneido
(integer)
```

Sets the seed for [Random](<../Random.md>). This is only required if you want a reproducible sequence of Random values.

Note: Each thread has its own separate sequence. Random.Seed only affects the current thread.