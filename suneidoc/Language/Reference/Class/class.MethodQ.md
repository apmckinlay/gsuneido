<div style="float:right"><span class="builtin">Builtin</span></div>

#### class.Method?

``` suneido
(string) => true or false
```

Returns true if the class or instance has the named method or false if not.

For example:

``` suneido
BorderControl.Method?("Resize")
    => true

BorderControl.Method?("Enlarge")
    => false
```

See also: [class.MethodClass](<class.MethodClass.md>)