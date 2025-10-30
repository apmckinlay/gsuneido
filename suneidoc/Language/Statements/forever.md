### forever
<pre>
<b>forever</b>
    <i>statement</i>
</pre>

Repeatedly executes the *statement*.

The **forever** statement is a shorthand for:

``` suneido
while true
    statement
```

or:

``` suneido
for (;;)
    statement
```

For example:

``` suneido
forever
    if x < 100
        x *= 2
    else
        break
```