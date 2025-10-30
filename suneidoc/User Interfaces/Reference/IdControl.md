### IdControl

``` suneido
(query, field, columns = false, fillin = false, from = false,
    mandatory = false, allowOther = false, whereField = false, width = 10, 
    readonly = false, noClear = false, restrictions = false, nameField = false, 
    abbrevField = false, access = false, prefixColumn = false, keys = false,
    restrictionsValid = true, filterOnEmpty = true, invalidRestrictions = false)
```

This is a single-line edit field with a list button on the right side. This 
control is used to provide a lookup for a value that exists in another table. 
For example, IdControl could be used for the customer field on invoices and it 
would provide a lookup on the customers table.

**nameField** and **abbrevField** are used to specify the fields from the query that are to be used as the display name and abbreviation matching.  If an abbreviation field is used, then only the abbreviation needs to be entered into the field and the name will be displayed. Also, if a prefix of either the name or abbreviation is entered and it is enough to uniquely identify the record, then the control will fill in all necessary fields from the record and that is all that must be entered.

NOTE: Although the nameField and abbrevField can be specified, it is recommended that a naming scheme be followed when naming these types of fields.  The field specified in the **field** argument should end in "_num".  The nameField and abbrevField field names should end in "_name" and "_abbrev" respectively.  This way the **nameField** and **abbrevField** arguments do not need to be specified because the control will use the above naming scheme to figure them out. Also, the IdFormat can then be used for the field.

The **restrictionsValid** parameter is used to distinguish whether or not values that are in the original query, but are omitted from the list because of restrictions, are still valid if the user manually enters them.

The **filterOnEmpty** argument is used to specify whether a query where the whereField's value is "" will return all the records in the table as if there was no where or only the records where the whereField's value is "". 

**invalidRestrictions** is used when you have restrictions on the query that are to be considered invalid by the control if the user types a value for a restricted record.

If the query needs to be restricted based on a field's value in the current record, then that field can be specified as the **whereField** argument. For example, you might have a contact lookup where you only want contacts for the current customer that is associated with the record. In this case the **whereField** would be the customer field. The field name from the current record must match the field name in the query. If this is not the case, then a rename can be done on the query to make the field names match.

**whereFieldKey** should be set to true if the **whereField** restriction is needed to make the lookup on the record unique.

To set up a simple lookup on invoices for the customer_num field, the customer_num field should have the following control specification: 

``` suneido
#(Id "customers" field: "customer_num")
```

To fillin the customer's address on the invoice:

``` suneido
#(Id "customers" field: "customer_num" fillin: (customer_address))
```

If the address field names don't match (customers table has customer_address1):

``` suneido
#(Id "customers" field: "customer_num"
    fillin: (customer_address) from: (customer_address1))
```

To only display customer name and address in the lookup list:

``` suneido
#(Id "customers" field: "customer_num" columns: (customer_name customer_address1))
```

For more information on the parameters, please refer to 
[KeyControl](<KeyControl.md>)

See also:
[IdFormat](<../../Reports/Reference/IdFormat.md>)