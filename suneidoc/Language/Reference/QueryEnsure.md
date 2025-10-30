### QueryEnsure

``` suneido
(table, rec)
```

Ensure that the table contains the given record.

If the record exists but is different it will be deleted.

If the record doesn't exist (or was different) it will be output.

See also: [QueryOutput](<../../Database/Reference/QueryOutput.md>)