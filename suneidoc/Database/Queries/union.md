### union

*query* **union** *query*

The result of union is a table with the same columns as both queries 
and with all the rows belonging to either query or both (with no duplicate rows in the result).  

**Note:** Both queries must have the same columns.

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

union

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
| 1 | 
| 2 | 
| 3 | 
| 4 | 
| 5 | 
| 6 | 
| 7 |

</div>
</div>