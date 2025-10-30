### while
<pre>
<b>while</b> <i>logical-expression</i>
    <i>statement</i>
</pre>

Repeatedly executes the *statement* as long as the *logical-expression* is true. The *statement* will not be executed at all if the *logical-expression* is initially false.

For example:

``` suneido
i = 0
while i < 10
    i += 3
```

**Note:** If logical-expression evaluates to something other than true or false, an exception will result: "conditionals require true or false".

See also: [do while](<do while.md>)