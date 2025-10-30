### join (natural join)

*query* **join** [ by (*fields*) ] *query*

The result of join is a table with all the columns from the input queries (without duplicates) and with the set of rows formed by combining each pair of rows with equal common columns. The input queries must have at least one column in common.

Notice that join does not output rows from the first table that do not have a matching row in the second table. If you need this, use [leftjoin](<leftjoin.md>) (left outer join).

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| table | tablename | 
| :---: | :---: |
| 17 | suppliers | 
| 18 | empty | 
| 19 | parts |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

join

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| table | column | 
| :---: | :---: |
| 17 | name | 
| 17 | city | 
| 19 | item | 
| 19 | cost |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| table | tablename | column | 
| :---: | :---: | :---: |
| 17 | suppliers | name | 
| 17 | suppliers | city | 
| 19 | parts | item | 
| 19 | parts | cost |

</div>
</div>

To control which fields to join on, use [rename](<rename.md>), [project](<project.md>), or [remove](<remove.md>). For example, if your orders table had a shipper1 and a shipper2 and you wanted to join by shipper1 with the shippers table (which has a shipper column):

``` suneido
orders rename shipper1 to shipper join shippers
```

Note: **by (*fields*)** is only an assertion, it does not alter which fields to join on. It only checks that the fields that the join uses are what you expect. If they differ an exception will be thrown. This is useful to ensure that the join does not change as the schema changes (e.g. you add fields).

**Note:** Only one-to-one joins are updatable.