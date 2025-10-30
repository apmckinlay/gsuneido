### IdFormat

``` suneido
(data = false, w = false, width = false, justify = "left",
    font = false, query = "", numField = "", nameField = false,
    abbrevField = false, abbrev = false)
```

Prints the name or abbreviation for an id value.

The **query** argument should be the query that the id value comes from.

**numField** should be the field from the **query** that contains the 
id value.

The **nameField** argument is used when the "_num, _name, _abbrev" naming convention is not being used.

If **abbrev** is true then the abbreviation field for the id will be printed instead of the name field.

The **abbrevField** argument is used to specify an abbreviation field in the query to match values against if abbrev is true.

NOTE: IdFormat assumes that the "_num, _name, _abbrev" naming convention is used for field names in the supplied **query** if no values are passed in as the nameField or abbrevField. For more information on this naming convention, see [IdControl.](<../../User Interfaces/Reference/IdControl.md>)

Example:

If you are printing customer_num on invoices, and you want the name of the 
customer to be printed, the format specification should look like:

``` suneido
   #(Id customers numField: "customer_num")
```

This example assumes there is a customers table set up with a customer_num 
field and a customer_name field.

See also: 
[IdControl](<../../User Interfaces/Reference/IdControl.md>)