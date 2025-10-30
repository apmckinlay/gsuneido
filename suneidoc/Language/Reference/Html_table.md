### Html_table

``` suneido
(object) => string
```

Generate an HTML table from an object.

The object should contain a list of objects, one for each row in the table.
The row objects should contain a list of values for the cells.
Both the main top-level object and the row objects may contain attributes
as named members.
The cell values can be plain numbers or strings
or they can be objects with the value as the first member
and other named members for attributes.

For example:

``` suneido
Html_table(Object(
    Object(#('Name', align: center), #('Age', align: center)),
    Object('Fred Smith', 56),
    Object('Sally Rogers', 35),
    border: 1, width: '50%'))
```

will return an HTML string that produces:

<div class="table-style table-half-width">

| Name | Age | 
| :---: | :---: |
| Fred Smith | 56 | 
| Sally Rogers | 35 | 

</div>

See also:
[Html_table_query](<Html_table_query.md>),
[Xml](<Xml.md>)