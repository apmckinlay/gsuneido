### times (product)

*query* **times** *query*

The result of times is a table with all the columns from both queries 
and with the set of rows formed by combining each possible pair of rows from the queries.  

**Note:** The queries cannot have any columns in common.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| a | b | 
| :---: | :---: |
| one | 12 | 
| two | 34 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

times

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| c | d | 
| :---: | :---: |
| three | 56 | 
| four | 78 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| a | b | c | d | 
| :---: | :---: | :---: | :---: |
| one | 12 | three | 56 | 
| one | 12 | four | 78 | 
| two | 34 | three | 56 | 
| two | 34 | four | 78 |

</div>
</div>