## QueryView

You can access the database interactively via QueryView, available from the
WorkSpace IDE menu. QueryView has two panes, a top Scintilla pane where you
type in database requests, and a bottom ListView pane that displays the
results of queries. Like the WorkSpace, F9/Run executes the current selection,
or if there is no selection, the non-blank lines around the insertion point.
![](<../res/queryview.png>)
The most common type of database request is a query, and the simplest kind
of query is just a table name.  Try running:

``` suneido
tables
```

This will show you the contents of the system table that contains
information about all the tables in the database.  Similarly:

``` suneido
columns
```

will show you information about all the columns in all the tables in the
database.  Usually you only want to see the columns for a single table.  But
wait, where is the table name?  Columns are related to tables by a table
number.  You could query tables to get the table number and then use that to
get the columns, but you might as well get Suneido to do that for you:

``` suneido
tables where tablename = "stdlib" join columns
```