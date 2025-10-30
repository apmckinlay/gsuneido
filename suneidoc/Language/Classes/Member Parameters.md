### Member Parameters

This is a shortcut to set members from the parameter values.

``` suneido
fn(.name) { ...
```

is equivalent to:

``` suneido
fn(name) { .name = name ...
```

and

``` suneido
fn(.Name)
```

is equivalent to:

``` suneido
fn(name) { .Name = name
```

For .Name the actual parameter is still "name"

Member parameters are commonly use for New methods (constructors).

One reason this is useful is because an explicit super call must be the first statement in a New method so it is not possible to do your own explicit member assignments prior to a super call.

See also:
[Use Member Parameters for Controls](<../../Appendix/Idioms/Use Member Parameters for Controls.md>)