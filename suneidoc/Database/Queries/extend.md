### extend

*query* **extend** *column* [ = *expression* ] [ , ... ]

The result of extend is a table with all the columns from the input plus the specified new columns and with all the rows from the input with values in the new columns calculated by evaluating the expressions (or rules) on each row.

Expressions are evaluated left to right, and any extended fields they depend on must come first. For example, `extend a = 1, b = 2 * a` is ok, but `extend b = 2 * a, a = 1` is illegal.

If no expression is given, the value can be calculated by a rule named: <code>Rule_ $ <i>column</i></code>

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| tablename | nrows | totalsize | 
| :---: | :---: | :---: |
| one | 100 | 5000 | 
| two | 2 | 700 |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

extend rowsize = totalsize / nrows

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| tablename | nrows | totalsize | rowsize | 
| :---: | :---: | :---: | :---: |
| one | 100 | 5000 | 50 | 
| two | 2 | 700 | 350 |

</div>
</div>

**Note:** If you project away fields required by an extend rule, the rule may not work properly. For more information see [project](<project.md>)