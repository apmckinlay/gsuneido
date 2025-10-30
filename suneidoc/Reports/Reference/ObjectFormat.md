### ObjectFormat

A [QueryFormat](<QueryFormat.md>) that gets its data from a list of records instead of a query.
Where you would specify the query, provide an object containing a list of records.

You can optionally specify:order
: An object containing a list of fields specifying the sort order. This is required to do Before's and After's. NOTE: The data must already be in this order. Use 
[object.Sort!](<../../Language/Reference/Object/object.Sort!.md>) if necessary.

columns
: An object containing a list of fields. This is a default and is only needed/used if no output format is specified.

For example:

``` suneido
data = QueryAll('columns')
data.Sort!(By(#table))
data.order = #(table)
data.columns = #(table, column)
fmt = ObjectFormat
    {
    Before_table(data)
        { return ['Text', '--- Table ' $ data.table] }
    }
Params.On_Preview([fmt data])
```