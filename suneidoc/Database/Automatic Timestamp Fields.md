## Automatic Timestamp Fields

If a column's name ends with _TS (upper case) then it will automatically be set to a current [Timestamp](<Reference/Timestamp.md>) when a record is output or updated. Any previous contents of the field will be overwritten. Automatic timestamp fields are useful for maintaining a "last updated" field.

_TS fields are automatically used by [AccessControl](<../User Interfaces/Reference/AccessControl.md>) to check if a record has been modified since it was read. (Without _TS fields the actual contents of the records are compared which can be awkward in some cases e.g. with rule fields.)

Since they would have identical values, only a single _TS field is supported. If there are multiple _TS columns only one of them will be updated.

For example:

``` suneido
Database("create tmp (a, b_TS) key(a)")
QueryOutput("tmp", Record(a: 1))
Query1("tmp where a = 1")

    => [a: 1, b_TS: #20060225.150657984]

QueryDo("update tmp set a = 2")
Query1("tmp where a = 2")

    => [a: 2, b_TS: #20060225.150759968]
```