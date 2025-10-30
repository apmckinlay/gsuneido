<div style="float:right"><span class="builtin">Builtin</span></div>

#### scanner.Next2

``` suneido
() => string or scanner
```

Returns the type of the next token from the scanner, one of:

``` suneido
#ERROR
#IDENTIFIER
#NUMBER
#STRING
#WHITESPACE
#COMMENT
#NEWLINE // whitespace containing one or more newlines
```

or "" (empty string) if none of the above (e.g. an operator)

Returns the scanner itself if no more tokens.

Next2 is intended to eventually replace Next. It is more efficient because it doesn't need to return the text.