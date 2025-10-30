<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Eval2

``` suneido
(callable, @args) => object
```

The only difference between object.Eval2 and [object.Eval](<object.Eval.md>) is that Eval2 returns the result inside an object. If there is no result (e.g. a return with no value) then the result object will be empty.