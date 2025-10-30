### Html_table_query

``` suneido
(query) => string
```

Generate an HTML table from a query.

On the query, use *project* to select and order the columns,
and *sort* to order the rows.

For example:

``` suneido
Html_table_query('tables where table < 10')
```

will return an HTML string that produces:

<div class="table-style">

| table# | tablename | nextfield | nrows | totalsize | 
| :---- | :---- | :---- | :---- | :---- |
| 2 | tables | 5 | 51 | 2864 | 
| 4 | columns | 3 | 159 | 3381 | 
| 6 | indexes | 9 | 76 | 4672 | 
| 8 | triggers | 3 | 0 | 100 | 

</div>

The column headings will use the heading from the *data dictionary*
Field_ definitions when available.

See also:
[Html_table](<Html_table.md>),
[Xml](<Xml.md>)