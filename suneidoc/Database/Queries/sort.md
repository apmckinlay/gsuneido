### sort
<pre>
query <b>sort</b> [ <b>reverse</b> ] column [ , ... ]
</pre>

Sort the result of a query by the specified column(s).
**sort** must be the last operation in a query.

For example:

``` suneido
columns sort table, column

tables sort reverse totalsize
```

Note: Strictly speaking, **sort** is not a relational query operator
since *relations* are not ordered.