### project

*query* **project** *columns*

The result of project is a table with only the specified columns from the input with any resulting duplicate rows removed.

**Note:** The input must contain the project columns.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| table | column | 
| :---: | :---: |
| 17 | name | 
| 17 | city | 
| 19 | item | 
| 19 | cost |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

project table

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| table | 
| :---: |
| 17 | 
| 19 |

</div>
</div>

Note: If you [extend](<extend.md>) a rule and then project away fields that are required by the rule, the rule may no longer work properly. For example:

``` suneido
mytable extend myrule project mykey, myrule // myrule may not work
```

Because rules are evaluated "lazily" i.e. on demand, myrule may not be evaluated until after the project. If myrule references other fields in mytable, they will no longer be available.

However, if the project does not include a key and therefore is not unique, then rules will be required to be evaluated in order to eliminate duplicates and they will work properly. For example:

``` suneido
mytable extend myrule project myrule // will work
```

See also: [remove](<remove.md>)