### MoneyBagTotal
<pre>
(append_function, currency_column, total_format = "Total", 
    <i>column</i>: moneybag ...)
</pre>

Calls the specified append_function to add _output format specifications to the report to print the totals from the moneybags in the specified columns.  One currency will be printed per line.  Normally used in After methods of a QueryFormat.

For example:

``` suneido
Total: ( (revenue, currency), (expense, currency) )
After_account(data)
    {
    MoneyBagTotalFormat(.Method('Append'), 'currency',
        revenue: data.total_revenue, expenses: data.total_expenses)
    }
After(data)
    {
    MoneyBagTotalFormat(.Method('Append'), 'currency', 'GrandTotal',
        revenue: data.total_revenue, expenses: data.total_expenses)
    }
```

See also:
[MoneyBag](<../../Language/Reference/MoneyBag.md>),
[QueryFormat](<QueryFormat.md>)