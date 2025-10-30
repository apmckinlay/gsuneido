### ImportXML

``` suneido
(from_file, to_query, fields = false, header = true)
```

Derived from [Import](<Import.md>)

Class for importing information from an *eXtensible Markup Language* textfile.

For example:

``` suneido
ExportXML('tables', 'tables.xml')
Database( "CREATE importTablesTest ( table , tablename , nextfield , nrows , totalsize ) 
    KEY( table )" )
ImportXML('tables.xml', 'importTablesTest')
```

Will take the information from the tables.xml file and output it to the importTablesTest,
which will look similar to this:

<div class='table-style'>

| table | tablename | nextfield | nrows | totalsize | 
| :---- | :---- | :---- | :---- | :---- |
| 2 | "tables" | 5 | 45 | 4726 | 
| 4 | "columns" | 3 | 145 | 3020 | 
| 6 | "indexes" | 9 | 65 | 2842 | 
| 8 | "triggers" | 3 | 0 | 100 | 

</div>

Note: This function can only handle XML files
that are in the same format as those produced by
[ExportXML](<ExportXML.md>).
Any DTD is ignored.

See also: [ExportXML](<ExportXML.md>)