### ensure
<pre>
<b>ensure</b> <i>table </i>( <i>columns </i>)
    <b>key</b> (<i>columns</i>) 
    <b>index</b> [ <b>unique</b> ] (<i>columns</i>) [ <b>in</b> table [ ( columns ) ] ]
</pre>

Ensure that the specified table has at least the specified characteristics. (It may have more.)

For example:

``` suneido
ensure mytable (name, salary) key(name)
```