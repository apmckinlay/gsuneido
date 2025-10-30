### sview
<pre>b>sview</b> <i>viewname</i> = <i>query</i></pre>

Define a *session* view. A *session* is a client connection when running client-server.

Session views:

-	do not show up in the views table
-	are unique to a session - different sessions have different sview's
-	only last till the end of a session, once you exit they are gone


For example:

``` suneido
sview my_tasks = tasks where owner = "apm"
```

or:

``` suneido
sview all_staff = active_staff union old_staff
```

A session view name may be used just like a table name in queries.

Use [drop](<../Requests/drop.md>) to un-define a session view.

See also: [view](<../Requests/view.md>)