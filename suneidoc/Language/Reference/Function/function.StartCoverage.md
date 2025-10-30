<div style="float:right"><span class="builtin">Builtin</span></div>

#### function.StartCoverage

``` suneido
(count = false)
```

From this point on, coverage information will be collected, until [function.StopCoverage](<function.StopCoverage.md>) is called.

If count is true then it will track the number of times each statement is executed. Otherwise it will only track whether or not each statement is executed.

Note: [CoverageEnable](<../CoverageEnable.md>)(true) must be called before the function is compiled.