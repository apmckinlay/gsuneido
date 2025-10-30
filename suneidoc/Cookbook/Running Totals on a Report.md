## Running Totals on a Report

**Category:** Reports

**Problem**

You want a column on a report that is a running total of another column.

**Ingredients**

[QueryFormat](<../Reports/Reference/QueryFormat.md>), [RowFormat](<../Reports/Reference/RowFormat.md>), [NumberFormat](<../Reports/Reference/NumberFormat.md>)

**Recipe**

Enter this in a library (e.g. mylib) as My_RunningTotal:

``` suneido
#(Params
    title: "Running Total Report"
    QueryFormat
        {
        Query: "tables"
        Before()
            {
            .running_total = 0
            return False
            }
        BeforeOutput(data)
            {
            .running_total += data.totalsize
            data.running_total = .running_total
            return False
            }
        Output: (Row tablename totalsize (Number field: running_total
            mask: '###,###,###', heading: 'Running\nTotal'))
        })
```

You can then run it from the WorkSpace with:

``` suneido
Window(My_RunningTotal)
```

**Discussion**

This example uses "tables" - a system table that stores information about the tables themselves. One of its fields is "totalsize" - the amount of space used by each table. This is the field that
we are making a running total of.

The first step is to initialize the running total in Before. Notice that we have to return False since we don't want to print anything.

Then, in BeforeOutput, we update the running total and then add it to the "data" - the current record. Again we return False since we don't want to print anything.

Finally, we specify a RowFormat for our Output. "tablename" and "totalsize" are actual fields in the table that have Field_ definitions so we can simply specify their field names. Since we haven't
made a Field_ definition for "running_total" we need to specify it's format (a NumberFormat). Notice the newline (\n) used to make a multi-line column heading.

**See Also**

Thanks to Steve Heyns (Zippy) for originally bringing up this topic
([http://www.suneido.com/forum/topic.asp?TOPIC_ID=632](<http://www.suneido.com/forum/topic.asp?TOPIC_ID=632>))

QueryFormat also has a built in "Total" facility for accumulating totals correctly for different sort breaks.