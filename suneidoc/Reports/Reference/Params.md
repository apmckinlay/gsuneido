### Params

``` suneido
(@report)
```

Displays a window with Print, PDF, Preview, and Page Setup buttons. If Print or Preview are chosen, a Report is created and used to output the supplied format to a printer or a PreviewControl.

Params looks for the following named arguments:
title
: Displayed at the top of the Params layout using a 
[TitleControl](<../../User Interfaces/Reference/TitleControl.md>).

name
: If the report has a name, the parameter data will be stored in the params table under that name, for the current user.

Params
: This is a user interface control specification (e.g. using FormControl) that allows the user to enter *parameters* for the report. These values can then be accessed from the report to place restrictions on queries, or change the way the report prints.

printParams
: Can be specified in order to print parameters in the page header. This must be an object containing the field names of the parameters to print.

SetParams
: If specified, SetParams should be an object where the members are field names and the values are what the fields are to be set to.  These values override any saved parameter values.

NoSaveLoadParams
: Set this to true so that saved params will not be loaded.  They will also not be saved.

validField
: Can be used to specify then name of a rule to be used for validating the Params. This rule should return "" if the parameters are valid and some message for the user if the parameters are invalid.

noPageRange
: If true, the print dialog will not allow the user to select a page range. This is useful when the report is performing some process where it is important to process all the way to the end.

The Page Setup is saved per computer (not per user) in a table called "devmode".