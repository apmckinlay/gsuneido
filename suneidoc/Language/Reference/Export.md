### Export

``` suneido
(from_query, to_file = false, fields = false, header = true)
```

Abstract base class for *exporting*

-	defaults the text file name to the table name
-	handles opening and closing the textfile
-	iterates through the query
-	returns the number of records exported


Classes derived from Export must implement:
Export1(x)
: Outputs the information from the record passed in, to the file.

Classes derived from Export have the option to implement:
Ext
: Sets a default file extension to use.

Before()
: Outputs information that must appear at the beginning of the file.

Header()
: Outputs information to the file that must appear first in the file or indicates that a Document Type Definition is also to be exported.

After()
: Outputs information that must appear at the end of the file.

Export provides:
From_query
: The query or table you are exporting data from.

To_file
: The file you are exporting the data to.  Defaults to the From_query with the Ext.

Tf
: The open file.

Fields
: The field names from the query or table.  Defaults to the column names from the query or table.

Header?
: Whether to output the header.  Defaults to true.

Export()
: Calls the Header() function if requested.
: Passes each record from the query or table to the Export1() function.
: Calls the Footer() function.
: Returns the number of records exported

Putline(line)
: Writes the line to the file

Close()
: Closes the file.

See also: 
[ExportCSV](<ExportCSV.md>),
[ExportTab](<ExportTab.md>),
[ExportXML](<ExportXML.md>),
[Import](<Import.md>)