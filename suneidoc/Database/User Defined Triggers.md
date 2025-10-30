## User Defined Triggers

Whenever records are output, updated, or deleted from a table, the system checks for a user defined trigger
with the name "Trigger_" $ tablename

Triggers take three arguments, for example:

``` suneido
Trigger_customer
function (transaction, oldrecord, newrecord)
    {
    }
```
transaction
: The containing transaction.  
This will be an update transaction since triggers are only called for update operations.

oldrecord
: For updates this is the version of the record prior to the update. 
For deletes this is the record being deleted.
For outputs this will be false.

newrecord
: For updates this is the new version. 
For deletes this will be false.
For outputs this will be the new record.

Triggers are only called **after** the operation succeeds.  
i.e. If the output, update, or delete fails, the trigger will not be called.

Any return value will be ignored.