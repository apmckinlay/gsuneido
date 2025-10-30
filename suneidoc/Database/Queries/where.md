### where (restriction/selection)

*query* **where** *expression*

The result of where is a table with the same columns as the input
but with just the rows where the expression is true.

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

where part = "abc"

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| date | part | quantity | 
| :---: | :---: | :---: |
| 00/09/01 | abc | 5 | 
| 00/09/02 | abc | 1 |

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

where quantity > 10 and   
      date = #20000902

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| date | part | quantity | 
| :---: | :---: | :---: |
| 00/09/02 | ghi | 20 |

</div>
</div>

Note: Query expressions do not support the full language expression syntax. In particular, calling member functions is not supported, only simple global functions. You can usually work around this by defining global helper functions to make any necessary method calls.