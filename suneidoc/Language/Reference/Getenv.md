<div style="float:right"><span class="builtin">Builtin</span></div>

### Getenv

``` suneido
( string ) => string
```

Return the value of an environment variable, or "" if the variable is not found.

For example:

``` suneido
Getenv("PATH")
    =>  "C:\\WINDOWS\\system32;C:\\WINDOWS"
```