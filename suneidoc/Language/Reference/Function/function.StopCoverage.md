<div style="float:right"><span class="builtin">Builtin</span></div>

#### function.StopCoverage

``` suneido
() => object
```

StopCoverage stops collecting coverage information and returns the accumulated results (since [function.StartCoverage](<function.StartCoverage.md>) was called).

The returned object will have source code positions as the member names. The value for each will either be a count if counting, or else true/false. There will be an entry for each tracked statement, even if it's zero or false.

Note: Constant propagation may mean that some statements do not generate any code and therefore are never tracked or executed. For example:

``` suneido
greet = "hello"
return greet
```

will be compiled as: `return "hello"` and `greet = "hello"` will never be executed.