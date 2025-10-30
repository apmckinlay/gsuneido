### QueryFormat

QueryFormat's are specified by the following members:

|  |  |  | 
| :---- | :---- | :---- |
| **Member** | **Type** | **Default** | 
| `Query` | string | none - must be supplied | 
| `Before` | format or false | false | 
| `Before_` | format or false | false | 
| <code>Before_<i>field</i></code> | format or false | false | 
| `BeforeOutput` | format or false | false | 
| `Output` | format or false | Row with Query columns | 
| `OutputFormat` | format or false | Output | 
| `AfterOutput` | format or false | false | 
| `After_` | format or false | false | 
| <code>After_<i>field</i></code> | format or false | false | 
| `AfterAll` | format or false | false | 
| `After` | format or false | false | 
| `Header` | format or false | Output.Header | 
| `Total` | object | `#()` | 
| `Count` | object | `#()` | 


These members may either be data members or methods returning the required data.
The only required member is Query, all other members are optional.

QueryFormat New copies its arguments into members (capitalized), so simple formats (with just data members) can be written without deriving a sub-class. For example:

``` suneido
Params(#(Query Query: 'tables'))
```

If there is an initial un-named argument, it is assumed to be the query, so the previous example can be shortened to:

``` suneido
Params(#(Query 'tables'))
```

If you need to specify methods for some of the members, derive a sub-class from QueryFormat.  For example:

``` suneido
Params(
    QueryFormat
        {
        Sort: 'sort tablename'
        Query()
            {
            return "tables " $ .Sort
            }
        }
    )
```

Arguments can be passed as members to derived formats. For example:

``` suneido
sort = "sort tablename"
Params(
    Object(
        QueryFormat
            {
            Query()
                {
                return "tables " $ .Sort
                }
            }
        Sort: sort
        )
    )
```

#### Output

Each row of the query result is generated using the Output format. Output defaults to a RowFormat with all the columns from the Query. So to specify the columns, you can either supply your own Output:

``` suneido
Params(#(Query 'tables', Output: (Row table, nrows)))
```

or add a project to the query:

``` suneido
Params(#(Query 'tables project table, nrows'))
```

Note: When Output is a method, it is only called **once** during setup; it is not called for each row.

Formats starting with an underscore, are forwarded to the Output. (e.g. _output sent to Row.Output).

#### OutputFormat

Specify OutputFormat instead of Output if you want to get column headings and use _output, but you do **not** want to generate the default output rows.

#### Before's and After's

`Before`, `After`, `BeforeOutput`, `AfterOutput`, `Before`_*field* and `After`_*field* are optional. They can be either formats, or methods that return formats or false. If they are methods they can accept an optional "data" argument. Before's and After's are grouped in a Vert with the row that triggered them, so they will not be split across a page break.

`Before`_*field* and `After`_*field* are only called if they are contained in the sort order of the query. Before's are called in the same order as the sort, and After's are called in the reverse order. For example, if the query sort order was `state, city` then the starting sequence would be:

``` suneido
Before
Before_state
Before_city
```

and the ending sequence would be:

``` suneido
After_city
After_state
After
```

Before's and After's are not called if the query is empty. If you want to produce some output even if the query is empty you can use the `AfterAll` method.

#### Header

`Query`'s Header simply delegates to the output format's `Header`, for `Row` this is the column headings (implemented via `ColHeadsFormat`).

#### Total and Count

Fields specified by `Total` and `Count` will be automatically totalled/counted. Totals are in the data object as "total_" $ field, e.g. "total_sales".  Counts are in the data object as "count_" $ field, e.g. "count_sales". Totaling/counting is done for each grouping determined by the query sort order. For example, if the sort order was `state, city` then in After_state the totals/counts would be for the state, in After_city the totals/counts would be for the city, and in After the totals/count would be for all the rows.

For example:

``` suneido
Params(#(Query
    Query: 'indexes sort table'
    Total: ('nnodes')
    After_table: (_output nnodes: (Total total_nnodes))
    After: (_output nnodes: (Total total_nnodes))
    ))
```

If you specify an object containing a pair of field names instead of just single field names, the first field name is taken as the field name of the amount, and the second is taken as the field name of the currency, and a [MoneyBag](<../../Language/Reference/MoneyBag.md>) is used to maintain the total.

``` suneido
Params(QueryFormat 
    {
    Query: 'transactions' 
    Total: ((amount, currency)),
    Output: ('Row', 'date', 'amount', 'currency')
    After(data)
        {
        fmt = Object('Vert', 'Hline')
        amounts = data.total_amount.Amounts()
        for (cur in amounts.Members())
            fmt.Add(Object('Row', 'date', 'amount', 'currency',
                data: Record(currency: cur, amount: amounts[cur])))
        return fmt
        }
    })
```

Notice that we had to sub-class QueryFormat in order to define After as a method.