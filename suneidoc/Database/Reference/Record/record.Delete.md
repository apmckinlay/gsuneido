<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Delete

``` suneido
()
```

Delete the record from the database. With no argument, the delete is done within the same transaction that read the record. (This transaction must not have been completed or rolled back.)

The record must have been read from the database.

**Note:** if an argument is supplied, this will be interpreted as [object.Delete](<../../../Language/Reference/Object/object.Delete.md>)