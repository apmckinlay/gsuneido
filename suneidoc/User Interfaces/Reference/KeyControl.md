### KeyControl

``` suneido
(query, field, columns = false, fillin = false, from = false,
    mandatory = false, allowOther = false, abbrevField = false,
    whereField = false, width = 10, readonly = false, noClear = false,
    restrictions = false, nameField = false, access = false, prefixColumn = false, 
    keys = false, filterOnEmpty = true)
```

This is a single-line edit field with a list button on the right side. This control is used to provide a lookup for a value that exists in another table. For example, KeyControl could be used for the customer field on invoices and it would provide a lookup on the customers table.

The **query** argument is used to specify the query to look the value up in. This can be as simple as just the table name, i.e. "customers".

The **field** argument is used to specify which field in the supplied query to look for the value in.

If **columns** are specified, then these are the columns that are displayed when the user brings up the list to choose from. If **columns** is not specified, then all columns resulting from the **query** are used.

The **fillin** argument is used if additional fields on the current record are to be filled in from the chosen id record.  It is specified as an object containing a list of fields. For example, you might have an address field on the customer record that you want filled in on invoice records when a customer is entered on the invoice.  If the field names from the id record do not match the fields on the current record, then the **from** argument can be used to match the fields up. It should be specified the same way the **fillin** argument is specified.

**mandatory** can be used when the user must enter a value in the field.

**allowOther** is used when the field can contain a value other than the values existing in the query.

If the query needs to be restricted based on a field's value in the current record, then that field can be specified as the **whereField** argument. For example, you might have a contact lookup where you only want contacts for the current customer that is associated with the record. In this case the **whereField** would be the customer field. The field name from the current record must match the field name in the query. If this is not the case, then a rename can be done on the query to make the field names match.

The **filterOnEmpty** argument is used to specify whether a query where the whereField's value is "" will return all the records in the table as if there was no where or only the records where the whereField's value is "".

The **width** parameter is used to set the width of the field.

If the field is to be protected, then the **readonly** argument should be passed as true.

If the **noClear** parameter is true, the fillin fields will not be cleared and filled with an empty object if there is no record to fill them with.

The **restrictions** argument is specified as a where clause (without the "where") and is used to restrict the query when entering new values, but they aren't used when setting the control to an existing value (restrictions aren't applied to old records).

The **abbrevField** argument is used to specify an abbreviation field in the query to match values against.  If the value entered exists in the query in the specified abbrevField, then that record will be used as a valid match.  KeyControl will also attempt to match prefixes entered to the values in the abbrevField in the query.

The **access** argument is used to add an 'Access' button on the list dialog. The value of the access argument should be the name of a control (often an AccessControl). e.g. `access: "MyAccess"` When the user clicks on the Access button on the list, the specified control will be opened in a dialog. This is normally used so the user can "access" the table that is the source for the list so that they can add or modify items without having to exit from the current screen.

The **prefixColumn** argument is used to sort the list and is the column that will be queried for a match to the users entry.

The **keys** parameter is used to pass the fields that the list can be searched on.

Examples:

To set up a simple lookup on the tables table for the table field:

``` suneido
   #(Key "tables" field: "table")
```

To fillin the table name:

``` suneido
   #(Key "tables" field: "table" fillin: (tablename))
```

If the table name field names don't match:

``` suneido
   #(Key "tables" field: "table" fillin: (name)
        from: (tablename))
```

To only display table and tablename fields in the lookup list:

``` suneido
   #(Key "tables" field: "table" columns: (table tablename))
```

See also: 
[IdControl](<IdControl.md>)