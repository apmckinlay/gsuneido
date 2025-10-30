<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Split

``` suneido
(separator = "") => object
```

Splits the string at occurrences of the specified separator string.

As of 2024-11-07 if the separator is empty it will split individual characters.

For example:

``` suneido
"one,two,three".Split(",") => Object("one", "two", "three")
```

**Note:** A trailing separator is ignored.

``` suneido
"one,two,three,".Split(",") => Object("one", "two", "three")
```

See also:
[object.Join](<../Object/object.Join.md>)[string.SplitOnFirst](<string.SplitOnFirst>)