### Import

``` suneido
(from_file, to_query, fields = false, header = true)
```

Abstract base class for *importing*

-	handles opening and closing the textfile
-	iterates through the file
-	returns the number of records imported


Classes derived from Import must implement:
Import1(line) => record
: Converts the line to a record.

Classes derived from Import have the option to implement:
Before()
: Inputs information that appears at the beginning of the file.
Header()
: Inputs information that appears at the beginning of the file.

Import provides:
from_file
: The file you are importing data from.

to_query
: The table you are importing the data to.

Tf
: The open file.

Fields
: The field names from file.  Defaults to false.

Header?
: Whether to input the header.  Defaults to true.

Import()
: Calls the Header() function if needed.
: Passes each line from the file to the Import1() function.
: Passes each record to the Output() function.

Getline() => string
: Returns the next line in the file.

Output(x)
: Starts a new transaction every 100 records.
: Outputs the record to the table.

Close()
: Completes any transactions that need to be.
: Closes the file.

See also:
[ImportCSV](<ImportCSV.md>),
[ImportTab](<ImportTab.md>),
[ImportXML](<ImportXML.md>),
[Export](<Export.md>)