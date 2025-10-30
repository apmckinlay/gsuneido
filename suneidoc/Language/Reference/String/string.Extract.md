<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Extract

``` suneido
(pattern, part = 0/1) => string or false
```

Returns all or part of the string that matches the given pattern,
or false if the string doesn't match the pattern.
If no part is specified then it returns part 1 if there is one,
otherwise it returns part 0 i.e. the entire match.

For example:

``` suneido
"hello world".Extract(".....$") => "world"
"hello world".Extract("(\w+) \w+") => "hello"
"hello world".Extract("(hello|howdy) (\w+)", 2) => "world"
"hello world".Extract("goodbye") => false
```

If the pattern matches in more than one place in the string, only the first match is used.

See also:
[Regular Expressions](<../../Regular Expressions.md>),
[string.Match](<string.Match.md>),
[string.Replace](<string.Replace.md>), 
[string.ExtractAll](<string.ExtractAll.md>)