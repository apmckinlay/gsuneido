<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/AccessControl/Messages">Messages</a><a href="/suneidoc/User Interfaces/Reference/AccessControl/Methods">Methods</a></span></div>

### AccessControl

``` suneido
(query, *control*, title = query, select = #(), startLast = false, 
    startNew = false, validField = false, stickyFields = #(), 
    locate = false, protectField = false, newOptions = false, 
    dynamicTypes = false, excludeSelectFields = false,
    saveOnlyLinked = false, historyFields = false, nextNum = false)
```

Creates a control used to 'access' the given **query**'s data. Access provides methods for navigating and modifying the data.

If the arguments contain a **title** member it will be used as the heading, otherwise the query will be used.

The **control** argument can either be a single control or it can be omitted. If no control is specified, query.Fields( ) will be used to create a form with one field per line, all in one group so they will be alligned. For example: 

``` suneido
(Access datadict)
```

If a single control is specified then it will be used. For example: 

``` suneido
(Access datadict (Vert name phone))
```

The single control format is the most flexible. For example, to use tabs: 

``` suneido
AccessControl('datadict'
  #(Vert 
    field 
    Skip 
    (Tabs 
      (Form control Tab: Control)
      (Form format Tab: Format))))
```

The **select** argument can be used to set an initial restriction on the query. This must be a container object with one object per restriction. Each restriction object must have three members: field, operator and value. 

For example:

``` suneido
select: #((name, '<', 'b') #(age, '>', 45))
```

The following operators are allowed:

``` suneido
=~   matches
!~   does not match
=    equals
==   equals
!=   not equal to
<    less than
<=   less than or equal to
>    greater than
>=   greater than or equal to
```

If the **startLast** argument is true, the access will start on the last record as opposed to the first.

If the **startNew** argument is true, the access will start on a new record as opposed to the first.

The **validField** argument can be used to specify a rule to be used for validation. This field (rule) does not necessarily have to be on the screen. This rule should return "" if the record is valid and some message for the user if the record is invalid.

The **stickyFields** argument is specified as an object containing names of fields that are to be sticky. When you enter new records, sticky fields will default to the value you entered on the last new record.

The **locate** argument is used to configure the locate control on the access. It should be an object with named members for keys and columns. The keys are the keys that will be allowed for use with the locate, and the columns are the columns that will be used for the pop-up list that comes up when the user is selecting a key value.

The **protectField** argument is used to specify a field name that will be used only to determine what fields to protect or not to protect. This field is a rule that evaluates to true, false, or an object. If the value is true then all fields will be protected. If it is false then none of the fields will be protected. If the value is an object, only the fields in the object will be protected. If the value is an object and the value at member 0 is 'allbut', all the fields will be protected except those in the object. If the value is an object and contains the member allowDelete with a value of true, then the record can be deleted whether or not there are any fields protected. Without this member set to true, the record will not be able to be deleted if there are any fields protected.

For example, a protect rule that returns a list of fields to protect would look like this:

``` suneido
Rule_protectRule
  {
  return #(name:, description:, date:, comments:)
  }
```

Or if you wanted to protect all the fields but the name field, and you wanted the user to be able to delete a record even if it has protected fields on it, your rule would look like this:

``` suneido
Rule_protectRule
  {
  return #('allbut', allowDelete:, name:)
  }
```

In both cases, your Access code would look like this:

``` suneido
(Access datadict (Vert name phone description date comments) protectField: 'protectRule')
```

The **newOptions** argument is used to specify the order of labels on the menu button. It should be specified in the following format:
<pre>
newOptions: #(<i>type_label1</i>, <i>type_label2</i>, <i>type_label3</i>)
</pre>

The **dynamicTypes** argument is used when the view needs to be changed on the fly based on a type field. If this argument is specified, then the New button will be changed to a menu button containing all of the specified types, so that different views can be used to enter records, depending on the type chosen. A default view can be specified to be displayed when the screen is entered. You can also specify types that should not show up under the New menubutton. These arguments should be specified in the following format:
<pre>
dynamicTypes: Object(typefield:<i> type_field_name</i>
  <i>type_label1</i>: #(value: <i>value</i> control: #(<i>control specifications</i>))
  <i>type_label2</i>: #(value: <i>value</i> control: #(<i>control specifications</i>))
  <i>default</i>: 'type_label1'
  <i>omitTypes</i>: #('type_label2')
  ...)
</pre>

Type_label1 and type_label2 would be replaced with the labels for the type. If the labels contain spaces then they must be enclosed in quotation marks. These labels will show up under the New menubutton and will also serve as the title of the Access when on records of the corresponding type. The actual values for each type that get saved in the typefield are specified in the value member of each type object.

The **excludeSelectFields** argument can be used to exclude fields from the Select option.

For example, if you did not want to display the group field in a select on the stdlib table:

``` suneido
AccessControl('stdlib', excludeSelectFields: #(group))
```

The **historyFields** argument is specified as an object containing names of fields that are to be shown on the History dialog under Current > History menu button. This argument should be specified in the following format:
<pre>
historyFields: #(createDate: <i>createDate_field_name</i>, 
    createBy: <i>createBy_field_name</i>,
    modifiedDate: <i>modifiedDate_field_name</i>, 
    modifiedBy: <i>modifiedBy_field_name</i>))
</pre>

The **nextNum** argument is an optional argument with the form:
<pre>
nextNum: #(field: <i>field_name_on_screen</i>,
    table: <i>table_name</i>, table_field: <i>field_name_in_table</i>)
</pre>

This will use [GetNextNum](<../../Database/Reference/GetNextNum.md>) to assign sequential numbers.

You can use GetNextNum.Create(table_name, table_field, starting_number) to create the table and an initial record in it.

See also:
[BrowseControl](<BrowseControl.md>)