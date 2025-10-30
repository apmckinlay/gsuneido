### create
<pre>
<b>create</b> <i>table </i>( <i>columns </i>)
    <b>key</b> [ <b>lower</b> ] (<i>columns</i>) 
    <b>index</b> [ <b>unique</b> ] (<i>columns</i>) 
        [ <b>in</b> table [ ( columns ) [ <b>cascade</b> [ <b>update</b> ] ] ]
</pre>

Create new tables. It will fail if the table already exists.

For example:

``` suneido
create mytable (name, city, age, Salary)
    key(name)
    index(city)
```

would create a table called "mytable" with the specified columns. The Salary column would be derived by calling Rule_salary. Two indexes will be created, one on name and one on city. The key index on name will ensure that no two records have the same name.

Multiple key's, and index's can be specified. Keys and indexes can be composed of multiple fields.

The difference between a key and a unique index is that unique index are optional - multiple records can have no value, whereas for a key only one record can have no value.