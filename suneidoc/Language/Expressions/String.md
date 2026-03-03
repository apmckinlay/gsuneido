### String

Operators:

|     | | 
| --- | ---- |
| $   | concatenate, returns a string | 
| =~  | matches regular expression, returns true or false | 
| !~  | does not match regular expression, returns true or false | 


For example:

``` suneido
"hello" $ " " $ "world"
    => "hello world"
```

Concatenation ($) will automatically convert booleans and numbers to strings. Using a string operator with other types will throw an exception of "cannot convert {type} to string". Use [Display](<../Reference/Display.md>) to convert other values to strings.

See also:
[Basic Data Types - String](<../Basic Data Types/String.md>),
[Cat](<../Reference/Cat.md>), 
[Match](<../Reference/Match.md>), 
[NoMatch](<../Reference/NoMatch.md>),
[Subscript](<Subscript.md>)