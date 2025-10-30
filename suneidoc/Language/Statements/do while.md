### do while
<pre>
<b>do</b>
    <i>statement</i>
    <b>while</b> <i>logical-expression</i>
</pre>

Repeatedly executes the *statement* as long as the *logical-expression* is true. The *statement* will always be executed at least once.

For example:

``` suneido
i = 0
do
    i += 3
    while i < 10
```

**Note:** If logical-expression evaluates to something other than true or false, an exception will result: "conditionals require true or false".

See also: [while](<while.md>)