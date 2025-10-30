### update

Update the records from a query:
<pre>b>update</b> <i>query</i> <b>set</b> <i>column</i> <b>=</b> <i>expression </i>[ , ... ]</pre>

For example:

``` suneido
update inventory where category = "bolts" set price = price * 1.1
```

**Note:** Only certain queries allow updates.