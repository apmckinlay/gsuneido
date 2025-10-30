### Asup

``` suneido
(string) => string
```

Replaces embedded expressions with their value.

For example:

``` suneido
s = "Hello default"
Asup(s)
    => "Hello default"
```

Note: Asup does not handle $> within the expressions (e.g. in a quoted string)