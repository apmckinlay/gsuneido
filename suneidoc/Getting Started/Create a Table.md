## Create a Table

As well as querying the database, you can also use QueryView for
administrative requests.  Let's create a table called mycontacts:

``` suneido
create mycontacts (firstname, lastname, phone, fax, email, website)
    key(lastname, firstname)
```

You can verify that this worked by querying your table.  You won't have any
data, but you'll see your columns.

**Note:** Tables must have at least one *key* - a unique identifying column or columns.

You can use SchemaView (from the WorkSpace IDE menu) to view the
*schema* (i.e. fields, keys, and indexes) of tables in the database.