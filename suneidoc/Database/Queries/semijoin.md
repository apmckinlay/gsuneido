### semijoin

*query* **semijoin** [ by (*fields*) ] *query*

The result of semijoin is a subset of the left hand source
with only the rows that have a matching row in the second query
based on equal common columns.
The input queries must have at least one column in common.

Semijoin is like [join](<join.md>) but only outputs columns from the first table,
and like [intersect](<intersect.md>) but matches on specified columns rather than requiring all columns to match.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">
	
| table | tablename |
| ----- | --------- |
| 17    | suppliers |
| 18    | empty     |
| 19    | parts     |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

semijoin

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

| table | column |
| ----- | ------ |
| 17    | name   |
| 17    | city   |
| 19    | item   |
| 19    | cost   |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

| table | tablename |
| ----- | --------- |
| 17    | suppliers |
| 19    | parts     |

</div>
</div>

Notice that row 18 from the first table is not included because there is no matching row in the second table.

To control which fields to join on, use [rename](<rename.md>), [project](<project.md>), or [remove](<remove.md>). For example, if your orders table had a shipper1 and a shipper2 and you wanted to semijoin by shipper1 with the shippers table (which has a shipper column):

``` suneido
orders rename shipper1 to shipper semijoin shippers
```

Note: **by (*fields*)** is only an assertion, it does not alter which fields to join on. It only checks that the fields that the semijoin uses are what you expect. If they differ an exception will be thrown. This is useful to ensure that the semijoin does not change as the schema changes (e.g. you add fields).