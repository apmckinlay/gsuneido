### String

Strings literals can be written with double or single quotes.  Within one kind of quote it is ok to include the other kind.

Backquotes (`) can be used when you do <u>not</u> want backslashes (\\) to have special meaning. This is useful for regular expressions and Windows file paths so you don't have to double the backslashes.

**Note:** Unlike many languages, Suneido allows multi-line strings.

If a string is also a valid identifier it can be written as #name. This is commonly used when passing a member name e.g. `ob.GetDefault(#mem, 0)` or a method name e.g. `ob.Map(#Trim)`

The `$` and `$=` operators are used to concatenate strings.

Additional user defined methods can be added by defining methods in a class called "Strings". (The stdlib standard library already includes a Strings class.)

Strings may include:

-	`\xhh` for hex character values, where h is 0-9, a-f, or A-F
-	`\t` for tab, `\n` for newline, `\r` for carriage return
-	`\'` for single quote, `\"` for double quote
-	`\0` for a nul (zero) character
-	`\\` for a literal backslash


For example:

``` suneido
"" // empty string'

' " '

"\""

"1\tIntroduction\n"

'a multi-line
string'

`c:\tmp\notes`
```

See also: [Expressions - String](<../Expressions/String.md>)