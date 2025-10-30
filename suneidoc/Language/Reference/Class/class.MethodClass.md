<div style="float:right"><span class="builtin">Builtin</span></div>

#### class.MethodClass

``` suneido
( method ) => class or false
```

Returns the class where the method is located or false if the method is not found.

For example:

``` suneido
BorderControl.MethodClass("Resize") => BorderControl
BorderControl.MethodClass("Construct") => Control
```

See also: [class.Method?](<class.Method?.md>)