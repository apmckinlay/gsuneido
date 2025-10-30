### GlobalPrint

``` suneido
(query)
```

Creates a report using Params.  Uses QueryFormat to output the query passed in.  
There are parameters for the report title, the columns to print, and the sort columns.

For example:

``` suneido
Window(GlobalPrint("tables"))
```