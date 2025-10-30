<div style="float:right"><span class="builtin">Builtin</span></div>

### CoverageEnable

``` suneido
(true or false)
```

Once coverage is enabled, any code compiled after this will include Cover instructions at the start of each statement. This adds some overhead to performance, in the range of 5 to 10% so it is not a good idea to enable this in production.

Coverage will not be actually tracked until [function.StartCoverage](<Function/function.StartCoverage.md>) or [class.StartCoverage](<Class/class.StartCoverage.md>) are called on particular functions.

Note: If the code was compiled prior to enabling coverage, you will need to force it to be recompiled using [Unload](<Unload.md>)