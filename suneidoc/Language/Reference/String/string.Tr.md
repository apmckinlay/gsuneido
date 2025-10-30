<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Tr

``` suneido
(from) => string
(from, to = "") => string
```

When passed a single argument, returns a new string with the characters in *from* deleted. For example, to delete blanks:

``` suneido
"one two three".Tr(' ') => "onetwothree"
```

If the from set begins with '^' then it is taken as all characters *except *the specified set.  For example, to delete everything but letters and blanks:

``` suneido
"one two three".Tr('^ ') => "  "
```

When passed two arguments, returns a new string with the characters in *from* replaced by the characters in *to*.

As with regular expressions, the shorthand a-z is expanded to all the characters from 'a' to 'z'. For example, to change upper case to lower case:

``` suneido
"HELLO World".Tr('A-Z', 'a-z') => "hello world"
```

If the to set is shorter than the from set, then the last character in the to set is replicated to make the sets the same length, and this replicated character is never put more than once in a row in the output.  For example, to change all runs of whitespace to single spaces:

``` suneido
"a     b\tc".Tr(" \t", " ") => "a b c"
```

See also:
[string.Replace](<string.Replace.md>)