### alter

Modify a table.  There are three variants:
<pre>b>alter</b> <i>table</i> <b>create</b> ...</pre>

Create new columns, keys, or indexes.  Will fail if the items already exist.

For example:

``` suneido
alter sales create (salesperson) index(salesperson)
```

See also: [create](<../Requests/create.md>)
<pre>b>alter</b> <i>table</i> <b>rename</b> <i>old_column_name</i> <b>to</b> <i>new_column_name</i> [ , ... ]</pre>

Rename one or more columns.  The data is not affected.

For example:

``` suneido
alter sales rename salesman to salesperson
```

See also: [rename](<../Requests/rename.md>)
<pre>b>alter</b> <i>table</i> <b>drop</b> ...</pre>

Delete columns, keys, or indexes.  Will fail if the items do not exist.

For example:

``` suneido
alter sales drop (salesperson)
```

**Note**: "delete" is an older alternative to "drop"

See also: [drop](<../Requests/drop.md>)