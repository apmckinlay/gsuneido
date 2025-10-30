### view
<pre>b>view</b> <i>viewname</i> = <i>query</i></pre>

Define a view.

For example:

``` suneido
view active_staff = staff where status = "active"
```

or:

``` suneido
view all_staff = active_staff union old_staff
```

A view name may be used just like a table name in queries.

Use [drop](<../Requests/drop.md>) to un-define a view.

See also: [sview](<../Requests/sview.md>)