### summarize

*query* **summarize** [ *by-columns*, ] [ *column* = ] *function column* [ , ... ]

The result of summarize is a table with the specified columns with one row for each value of the by-columns. The allowable functions are: max, min, total, average, count, and list. If a name is not specified for a calculated column, the function name plus the original column name will be used.

**Note:** If there are no rows in the input, there will be no output rows.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| date | part | quantity | 
| :---: | :---: | :---: |
| 00/09/01 | abc | 5 | 
| 00/09/01 | def | 10 | 
| 00/09/02 | abc | 1 | 
| 00/09/02 | ghi | 20 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

summarize date,    
      count, total quantity

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| date | count | total_quantity | 
| :---: | :---: | :---: |
| 00/09/01 | 2 | 15 | 
| 00/09/02 | 2 | 21 |

</div>
</div>

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| date | part | quantity | 
| :---: | :---: | :---: |
| 00/09/01 | abc | 5 | 
| 00/09/01 | def | 10 | 
| 00/09/02 | abc | 1 | 
| 00/09/02 | ghi | 20 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

summarize part,    
total quantity

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| part | total_quantity | 
| :---: | :---: |
| abc | 6 | 
| def | 10 | 
| ghi | 20 |

</div>
</div>
&nbsp;
<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| date | part | quantity | 
| :---: | :---: | :---: |
| 00/09/01 | abc | 5 | 
| 00/09/01 | def | 10 | 
| 00/09/02 | abc | 1 | 
| 00/09/02 | ghi | 20 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

summarize count,   
      total quantity

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| count | total_quantity | 
| :---: | :---: |
| 4 | 36 |

</div>
</div>

As a useful special case, when you get the overall minimum or maximum of a key field, you also get the associated record. For example:

``` suneido
Query1('tables summarize min table')
    => [min_table: 1, tablename: "tables", table: 1]
```