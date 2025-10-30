### ImportCSV

``` suneido
(from_file, to_query, fields = false, header = true)
```

Derived from [Import](<Import.md>)

Class for importing information from a *Comma Separated Value* textfile.

For example:

``` suneido
ExportCSV('tables', 'tables.txt')
Database( "CREATE importTablesTest ( table , tablename , nextfield , nrows , totalsize ) 
    KEY( table )" )
ImportCSV('tables.txt', 'importTablesTest')
```

Will take the information from the tables.txt file and output it to the importTablesTest table, 
which will look similar to this: 

<div class='table-style'>

| table | tablename | nextfield | nrows | totalsize | 
| :---- | :---- | :---- | :---- | :---- |
| 2 | "tables" | 5 | 45 | 4688 | 
| 4 | "columns" | 3 | 145 | 3020 | 
| 6 | "indexes" | 9 | 65 | 2842 | 
| 8 | "triggers" | 3 | 0 | 100 | 

</div>

**Note:** If header is set to false, the fields must be passed in.

See also: [ExportCSV](<ExportCSV.md>)