### QueryOutput

``` suneido
(query, record)
```

Outputs the record to the query in a standalone update transaction.

For example:

``` suneido
QueryOutput("phonelist", Record(name: "Sue", phone: "123-4567"))
```