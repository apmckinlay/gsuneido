## Report Control Breaks

When a sort is specified on a report's query, sort breaks can be defined to print summaries, totals, or do calculations before and after each value in the sort. There are also breaks that can be used before and after the report query.

Assuming you have a QueryFormat that has the query "sales sort salesperson, city", the following methods can be defined to handle the control breaks:

-	Before()
-	Before_salesperson(data)
-	Before_city(data)
-	BeforeOutput(data)
-	AfterOutput(data)
-	After_salesperson(data)
-	After_city(data)
-	After(data)
-	AfterAll()


The above can be defined as members or methods of QueryFormat, but must evaluate to a format specification or false.  The exception is Total, this must be an object containing the list of fields that the QueryFormat will total automatically.

The following is a simple example of how to use control breaks in a report.

From QueryView create a sales table with the following command:

``` suneido
    create sales (sale_id, salesperson, sale_amount, sale_city)
        key (sale_id)
```

In LibraryView create Field_ definitions for each field, for example:

**`Field_sale_amount`**
``` suneido
Field_dollars
    {
    }
```

Next, use IDE > Access a Query or IDE > Browse a Query to add some records to the table. Make sure you have a few different salespersons and cities.

Now we can define a report (in LibraryView) to print the sales using the sort breaks mentioned above.

**`Sales_Report`**
``` suneido
#(Params
    title: "Sales Report"
    name: "Sales_Report"
    QueryFormat
        {
        Query()
            {
            return "sales sort salesperson, sale_city"
            }
        Before()
            {
            return #("Text", "Sales Report by Salesperson, City")
            }
        Before_salesperson(data /*unused*/)
            {
            return #(Hline)
            }
        Before_sale_city(data)
            {
            return Object("Text", "Sales for: " $ data.salesperson $ ", City: " $ data.sale_city)
            }
        Total: (sale_amount)
        After_salesperson(data /*unused*/)
            {
            return #(_output
                sale_amount: (Total "total_sale_amount"))
            }
        After(data /*unused*/)
            {
            return #(_output
                sale_amount: (GrandTotal "total_sale_amount"))
            }
        AfterAll()
            {
            return #("Text", "END OF SALES REPORT")
            }
        }
    )
```