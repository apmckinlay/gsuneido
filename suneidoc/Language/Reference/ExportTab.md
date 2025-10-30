### ExportTab

``` suneido
(from_query, to_file = false, header = true)
```

Derived from [Export](<Export.md>)

Class for exporting information as a *Tab Delimited* textfile.

For example:

``` suneido
ExportTab('tables', 'tables.txt', header:) 
```

Will take the information out of the tables table and output it to the 
tables.txt file in a format similar to this:

``` suneido
table   tablename   nextfield   nrows   totalsize
2   tables  5   45  4688
4   columns 3   145 3020
6   indexes 9   65  2842
8   triggers    3   0   100
```

Specifying header: produces the initial line with the field names.

See also: [ImportTab](<ImportTab.md>)