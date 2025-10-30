### ExportCSV

``` suneido
(from_query, to_file = false, header = "Fields")
```

Derived from [Export](<Export.md>)

Class for exporting information as a *Comma Separated Value* textfile.

For example:

``` suneido
ExportCSV('tables', 'tables.txt', header:)
```

Will take the information out of the tables table and output it to the
tables.txt file in a format similar to this:

``` suneido
"table","tablename","nextfield","nrows","totalsize"
2,"tables",5,45,4688
4,"columns",3,145,3020
6,"indexes",9,65,2842
8,"triggers",3,0,100
```

The **header** argument can be one of the following:
`"Fields"`
: gives a header line with the field names (this is the default)

`"Prompts"`
: gives a header line with prompts using 
[Prompt](<Prompt.md>)(field)

`"None"`
: no header line, just the data

See also: [ImportCSV](<ImportCSV.md>)