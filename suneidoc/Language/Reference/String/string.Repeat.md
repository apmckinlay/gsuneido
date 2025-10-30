<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Repeat

``` suneido
(count) => string
```

Returns a string consisting of count copies of the string.

If count is less than or equal to zero, the result is "".

For example:

``` suneido
"hello".Repeat(2) => "hellohello"
```