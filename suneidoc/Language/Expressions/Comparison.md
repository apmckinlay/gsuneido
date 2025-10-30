### Comparison

Comparison operators always produce a boolean result (true or false).

|  |  | 
| :---- | :---- |
| `is` | equals | 
| `isnt` | not equals | 
| `<` | less than | 
| `<=` | less than or equal to | 
| `>` | greater than | 
| `>=` | greater than or equal to | 
| `=~` | matches regular expression | 
| `!~` | doesn't match regular expression | 
| `in` | equals one of a list of expressions | 
| `not in` | not equal to any of a list of expressions | 


**Note**: Strings and numbers are never equal e.g. 123 is not equal to "123". Use the String() or Number() functions to convert values.

**Note**: Object and Record equality do not apply defaults or rules.

`in` and `not in` are used with a single value on the left and a parenthesized list of values on the right. For example:

``` suneido
x in (n - 1, n, n + 1)
```

See also: [Regular Expressions](<../Regular Expressions.md>), [Gt](<../Reference/Gt.md>)