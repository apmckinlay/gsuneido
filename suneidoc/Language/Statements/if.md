### if
<pre>
<b>if</b> <i>logical-expression</i>
    <i>statement1</i>

<b>if</b> <i>logical-expression</i>
    <i>statement1</i>
<b>else</b>
    <i>statement2</i> 
</pre>

If the *logical-expression* is true, *statement1* is executed, if false and there is an else, *statement2* is executed.

For example:

``` suneido
if i < 0
    Print("negative")
else
    Print("positive")
```

**Note:** If logical-expression evaluates to something other than true or false, an exception will result: "conditionals require true or false".

See also: [Conditional (?:) Operator](<../Expressions/Conditional.md>)