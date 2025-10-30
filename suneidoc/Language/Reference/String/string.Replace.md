<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Replace
<pre>pattern, replacement, count = <i>all</i>) => string</pre>

Creates a new string with *pattern* replaced with *replacement*.  
The default is to replace all occurrences of the pattern, 
but a count can be specified to do only a certain number of replacements (e.g. only 1).
The replacement can be either a string or a callable value (e.g. a block or function)
that returns a string.

For example:

``` suneido
"hello world".Replace("o", "O") => "hellO wOrld"
"hello world".Replace("o", "O", 1) => "hellO world"
```

The replacement string can contain special character sequences that are expanded as follows:

<div class="table-full-width">

|  |  | 
| :---- | :---- |
| `&` or `\0` | the string that was matched by the pattern | 
| `\1, \2, ` ... | the portion of the string that was matched by the nth parenthesized part of the regular expression, counting opening parentheses from the left | 
| `\u` | convert the single following character to upper case | 
| `\l` | convert the single following character to lower case | 
| `\U` | convert all the following characters to upper case (until \E) | 
| `\L` | convert all the following characters to lower case (until \E) | 
| `\E` | end `\U` or `\L` | 
| `\\` | a literal '\' (backslash) | 

</div>

For example:

``` suneido
"Hello World".Replace("(\w+) (\w+)", "\L\2 \U\1")  =>  "world HELLO"
"hello world".Replace("\w+", "\u&")  =>  "Hello World"
```

If the replacement string starts with `\=`
this turns off the special meaning of the characters described above.

If the replacement is a callable value, 
it is passed the matched portion of the string
and returns the replacement string.
It is called once for each match.
For example:

``` suneido
"abc".Replace('.', { |s| s.Upper() } )
    => "ABC"
```

If the replacement callable value does not return anything,
then no replacement is done.
This can be useful for the side-effects.
For example:

``` suneido
"hello world".Replace('[a-z]+') { |s| Print(s) }
    hello
    world
    => "hello world"
```

See also:
[Regular Expressions](<../../Regular Expressions.md>),
[string.Extract](<string.Extract.md>),
[string.Match](<string.Match.md>),
[string.Tr](<string.Tr.md>)