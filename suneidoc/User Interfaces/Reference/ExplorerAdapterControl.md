### ExplorerAdapterControl

``` suneido
( control, field )
```

Adapts a single-field control, whose Get and Set work with single values, to be used as with Explorer, that requires Get and Set to work with objects.

``` suneido
Set(object) => control.Set(object[field])

Get() => Object(field: control.Get())
```

Forwards Dirty? to the control.

An example of the
[Adapter](<../../Appendix/Patterns/Adapter.md>)
pattern.