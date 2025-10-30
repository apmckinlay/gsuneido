### Conditional
<pre>i>logical-expression</i> ? <i>expression1</i> : <i>expression2</i></pre>

The logical expression is evaluated,
if it is true, then expression1 is evaluated and is the result,
if it is false, then expression2 is evaluated and is the result.
If logical-expression evaluates to something other than true or false,
an exception will result: "conditionals require true or false".

**Note:** Only one of *expression1* or *expression2*
will be evaluated.

For example:

``` suneido
if (condition)
    x = value1
else
    x = value2

if (condition)
    return value1
else
    return value2
```

can be written as:

``` suneido
x = condition ? value1 : value2

return condition ? value1 : value2
```