<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Writeline

``` suneido
(string) => string
```

Output a string to a file object followed by a newline and returns the string.

For example:

``` suneido
File("tmp", "w")
    {|f|
    f.Write("hello world")
    }
```

This would create a file called "tmp" containing:

``` suneido
hello world\n
```

See also:
[file.Write](<file.Write.md>)