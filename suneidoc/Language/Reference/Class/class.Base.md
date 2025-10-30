<div style="float:right"><span class="builtin">Builtin</span></div>

#### class.Base

``` suneido
() => class or false
```

Returns the base (parent) super class of a class or instance, or false for a class without a base class.

For example:

``` suneido
VertControl.Base() => Group
Stack.Base() => false
Stack().Base() => Stack // an instance's base is its class
```

See also: [class.Base?](<class.Base?.md>)