### RowHeaderFormat

``` suneido
( fields )
```

Prints centered column headings with a line underneath across the entire column width.

**fields** is a list of fields to print column headings for.

RowHeaderFormat can be used instead of [ColHeadsFormat](<ColHeadsFormat.md>) for standard column headings. Very useful for group reports and simple to use.

**PS: **If you don't define a QueryFormat or deactivate **Header** in defined [QueryFormat](<QueryFormat.md>), the column headings will still be printed on top of the page.

Example in a library:
<pre>
#(Params
    name: 'tables_Report'
    title: "Suneido Tables and Fields"
    //Use QueryFormat to build a composite report.
    QueryFormat
        {
        // Build the query.
        Query()
            {
            return "columns join tables remove nextfield
              sort tablename,field, nrows, totalsize"
            }
        // Prints only these fields on a row.
        Output()
            {
            // Control the size of field print with width.
            // Without width, the column is too small to fit the text.
            return #('Row' (column width: 30), field)
            }
        <b>// Deactivate the printing of column headings on top of page.
        Header()
            {
            return false
            }</b>
        // Print the tablename (Group Header).
        Before_tablename(data)
            {
            return Object('Vert'
                #(Text '')
                Object('Horz'
                    Object('Text', data.tablename, justify: 'left',
                      font: #(name: 'Arial' size: 14 weight: 1200))
                    #('Hfill')
                    Object('Text', data.tablename, justify: 'right',
                      font: #(name: 'Arial' size: 14 weight: 1200))
                    )
                #('Hline' thick: 30 before: 30 after: 1)
                <b>// Print the column headings.
                #('RowHeader' (column width: 30), field)</b>
                )
            }
        // Group Break.
        After_tablename(data)
            {
            return #('Vert'
                ('Text' '')
                )
            }
        // Print a line separator between records.
        AfterOutput()
            {
            return #(Hline thick: 1 before: 20 after: 1)
            }
        }
    )
</pre>

Same example in workspace:
<pre>
Params(
    QueryFormat
        {
        Query()
            {
            return "columns join tables remove nextfield
              sort tablename,field, nrows, totalsize"
            }
        Output()
            {
            return #('Row' (column width: 30), field)
            }
        Header()
            {
            return false
            }
        Before_tablename(data)
            {
            return Object('Vert'
                #(Text '')
                Object('Horz'
                    Object('Text', data.tablename, justify: 'left',
                      font: #(name: 'Arial' size: 14 weight: 1200))
                    #('Hfill')
                    Object('Text', data.tablename, justify: 'right',
                      font: #(name: 'Arial' size: 14 weight: 1200))
                    )
                #('Hline' thick: 30 before: 30 after: 1)
                <b>#('RowHeader' (column width: 30), field)</b>
                )
            }
        After_tablename(data)
            {
            return #('Vert'
                ('Text' '')
                )
            }
        AfterOutput()
            {
            return #(Hline thick: 1 before: 20 after: 1)
            }
        }
    name: 'tables_Report'
    title: "Suneido Tables and Fields"
    )
</pre>