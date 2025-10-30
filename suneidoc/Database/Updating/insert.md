### insert

Used to insert a record into a query:
<pre>b>insert</b> { <i>column</i>: <i>value </i>[ , ... ] } <b>into</b> <i>query</i></pre>

For example:

``` suneido
insert { name: "Joe", salary: 37500 } into employees
```

Or to insert records from a query into an existing table:
<pre>b>insert</b> <i>query</i> <b>into</b> <i>table</i></pre>

For example:

``` suneido
insert sales where city = "Fargo" into fargo_sales
```