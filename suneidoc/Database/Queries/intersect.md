### intersect

*query* **intersect** *query*

The result of intersect is a table with the same columns as the input 
and with all rows belonging to **both** input queries.  

**Note:** Both queries must have all the same columns.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

|             number | 
| :---: |
| 1 | 
| 2 | 
| 3 | 
| 5 | 
| 6 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

intersect

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

|             number | 
| :---: |
| 2 | 
| 3 | 
| 4 | 
| 6 | 
| 7 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

|             number | 
| :---: |
| 2 | 
| 3 | 
| 6 |

</div>
</div>