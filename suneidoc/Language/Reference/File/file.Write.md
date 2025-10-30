<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Write

``` suneido
(string) => string
```

Output a string to a file object and returns the string.

For example:

``` suneido
File("tmp", "w")
    {|f|
    f.Write("hello world")
    }
```

This would create a file called "tmp" containing:

``` suneido
hello world
```

See also:
[file.Writeline](<file.Writeline.md>)