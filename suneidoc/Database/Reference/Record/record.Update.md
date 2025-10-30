<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Update

``` suneido
() or (object)
```

Normal usage is to read a record, modify it, and then call its Update method to save the changes. For example:

``` suneido
Transaction(update:)
    { |t|
    x = t.Query1(...)
    x.abc = ...
    x.Update()
    }
```

If a record is passed in, it is used as the source of the new data (instead of the record you call Update on). This is useful when you read a record in one transaction, and then want to re-read it and update it in another transaction. For example:

``` suneido
x = Query1(...)
...
Transaction(update:)
    { |t|
    y = t.Query1(...)
    y.Update(x)
    }
```

However, be aware this could overwrite another user's changes in a multi-user situation.