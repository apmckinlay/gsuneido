<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/TreeModel/Methods">Methods</a></span></div>

### TreeModel

``` suneido
(table)
```

An ExplorerControl model.

Accesses a tree stored in a database table.  
The table must have the following specification (presumably with additional fields for data):

``` suneido
(num, parent, group, name, ...) key (num) index (parent)
```

Each record is assigned a unique number (starting at 1).  
In the table group = -1 for leaf items, group == parent for groups.  

A key of name,group can be used to ensure unique item names overall, 
and unique group names within groups.