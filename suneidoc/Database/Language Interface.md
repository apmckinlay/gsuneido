## Language Interface

### Requests

Requests are executed using the Database() function.  For example:

``` suneido
Database("create phonelist (name, phone) key(phone)")
```

### Queries and Updating

Queries and update requests are executed using the Query method of a transaction.
For example:

``` suneido
t = Transaction(update:)
t.Query("delete customers where city = 'Regina'")
t.Complete()
```

Notice that an update transaction is required if you are modifying the database.
Also notice the use of single quotes to avoid having to escape nested quotes.

Where applicable, it is often better to use the block form of Transaction 
to ensure that the transaction will be ended.  For example:

``` suneido
Transaction(update:)
    { |t|
    t.Query("delete customers where city = 'Regina'")
    }
```

Or use the QueryDo function in stdlib:

``` suneido
QueryDo("delete customers where city = 'Regina'")
```

Insert, update, and delete queries return the number of records affected.

For a query rather than an update request, transaction.Query returns a query object
which is then used to read records.  For example:

``` suneido
Transaction(read:)
    { |t|
    q = t.Query("customers where city = 'Saskatoon'")
    while (false isnt (x = q.Next()))
        Print(x)
    }
```

Notice that in this case we only needed a read transaction.

The QueryApply function in stdlib simplifies this to:

``` suneido
QueryApply("customers where city = 'Saskatoon'")
    { |x| Print(x) }
```

You can also use Query1 when your query should only match a single record
(i.e. you are selecting on a key, or the table has an empty key).
Or you can use QueryFirst and QueryLast.  For example:

``` suneido
record = Query1("tables where tablename = 'stdlib'")
```

Once you've read a record from a query, you can update it or delete it.  For example:

``` suneido
QueryApply("inventory where price < 100")
    { |x|
    if (x.price < .1)
        x.Delete()
    else
        {
        x.price = x.price * 1.1
        x.Update()
        }
    }
```

Records can also be output to a query
(as long as the query is one that permits outputs).
For example:

``` suneido
Transaction(update:)
    { |t|
    q = t.Query("customers")
    q.Output(record)
    }
```

The QueryOutput function in stdlib simplifies this to:

``` suneido
QueryOutput("customers", record)
```

Either Object's or Record's can be output.  Records support
[Rules](<Rules.md>).

### Examples

To export data to a tab delimited text file:

``` suneido
fields = #(name, phone, fax, email)
File("phonelist.txt", "w")
    { |f|
    QueryApply("phonelist")
        { |x|
        line = ""
        for (field in fields)
            line $= x[field] $ '\t'
        f.Writeline(line[..-1])
        }
    }
```

To import data from a tab delimited text file:

``` suneido
fields = #(name, phone, fax, email)
Database("ensure phonelist (name, phone, fax, email) key(name)")
File("phonelist.txt")
    { |f|
    while (false isnt (line = f.Readline()))
        {
        record = Record()
        values = line.Split('\t')
        for (i in values.Members())
            record[fields[i]] = values[i]
        QueryOutput("phonelist", record)
        }
    }
```