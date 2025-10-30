<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/BrowseControl/Methods">Methods</a></span></div>

### BrowseControl

``` suneido
(query, columns = false, linkField = false, stickyFields = false, validField = false, 
    protectField = false, title = false, statusBar = false, headerFields = #(), 
    dataMember = "browse_data", noShading = true, noHeaderButtons = false, notifyLast = false, 
    columnsSaveName = false, stretch = false, alwaysReadOnly = false)
```

Displays records from a query in list format and allows records to be added, modified and deleted.

**Note:** This form of browse reads all the records into memory - it should not be used on very large queries, could use [VirtualListControl](<VirtualListControl.md>) instead

The **query** argument is used to specify the query that the records will be taken from. 

Use the **columns** argument to specify which columns are to be displayed in the list.  This should be passed as an object containing field names from the query.

The **linkField** argument is used when you have line items on an Access screen and a BrowseControl is used to display the line items.  The link field should be a key from the header table.  This also requires that a foreign key exists between the header and line item tables.

The **title** argument is displayed in a large font just above the browse header. If no title is passed then the query is used.

The **statusBar** argument determines whether there will be a status bar or not.

If **linkField** is being used to link the browse items to a header record and the line item records need access to the data from the header record, then the **headerFields** argument can be specified as a list of fields and then these header fields will be added to each line item record.

Also, in the case where the browse is linked to a header record, the browse will automatically set a field in the header's data to its own data so that the header record can access the browse records.  This is useful in rules.  The field name on the header record defaults to "browse_data", however if there is more than one browse linked to a header record, then there will be a conflict with the header's "browse_data" field.  To rename the field being used by the header to store the browse data, pass the new field name as a string to the browse as the **dataMember** argument.

The **noShading** argument determines whether or not every second line in the browse will be shaded.

The **noHeaderButtons** argument determines whether or not header buttons can be clicked.

If the **notifyLast** argument is true, then the current browse will be saved last.  This is useful when you have multiple browses embedded in an Acess and would like to save a particular browse last.

The **columnsSaveName** parameter is used for cases where the browse does not have an initial query and there is no title to save the columns widths and order under. The column information will be saved under this name.

If **stretch** is true, ListStretchControl will be used instead of ListControl. This will automatically resize the last column to fill the window.

If **alwaysReadOnly** is true, the Browse will always be protected.

For example, to create a simple Browse on the tables table

``` suneido
BrowseControl("tables")
```

Would display:

![](<../../res/browse.png>)

For more information on stickyFields, validField, and protectField parameters, please refer to [AccessControl](<AccessControl.md>).

Note: BrowseControl delegates any unknown methods to its ListControl. So methods from [ListControl](<ListControl.md>) are also available.

See also:
[VirtualListControl](<VirtualListControl.md>)[Tools/Browse](<../../Tools/Browse.md>),
[ExplorerListViewControl](<ExplorerListViewControl.md>),
[ListViewControl](<ListViewControl.md>)