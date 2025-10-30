## System Tables

Information about the database (the schema) can be viewed in the following system tables:

tables - lists the tables that exist

columns - lists the columns for each table

indexes - lists the keys and indexes for each table

views - lists the view definitions

history - lists potential [transaction.Asof](<Reference/Transaction/transaction.Asof.md>) date-times. Viewing the history may be slow because it has to scan potentially the entire database file. And if the database has not been compacted recently, there may be a large number of entries.

The contents of these tables may only be altered by the system.  However, they can be read from just like any other table.

**Note**: These are not "physical" tables. They are virtual tables that are "views" of internal metadata. **Warning**: Accessing them may be slow.