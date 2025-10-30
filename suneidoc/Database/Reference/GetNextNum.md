### GetNextNum

Methods:
`Create(table, field = "nextnum", nextnum = 1)`
: Create the next number table.

`Reserve(table, field = "nextnum")`
: Returns the next number to be assigned. This only *reserves* the number, the *reservation* will expire in GetNextNum.ReserveSeconds unless it is Renew'ed

`Renew(num, table, field = "nextnum")`
: Renew a *reservation*. This must be done prior to GetNextNum.ReserveSeconds to keep the *reservation*.

`Confirm(num, table, field = "nextnum")`
: Makes a reservation permanent.

`PutBack(num, table, field = "nextnum")`
: Cancel a *reservation*, allowing a future Reserve to return the number.

`ChangeNextNum(table, field, nextnum)`
: 

Used by [AccessControl](<../../User Interfaces/Reference/AccessControl.md>) nextNum option.