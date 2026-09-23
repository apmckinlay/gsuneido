## System Tables

Information about the database (the schema) can be viewed in the following system tables.
In gSuneido these are not actual tables, they are virtual tables based on internal information.

tables - the tables that exist

columns - the columns for each table

indexes - keys and indexes for each table

views - the view definitions

history - potential [transaction.Asof](<Reference/Transaction/transaction.Asof.md>) date-times. Viewing the history may be slow because it has to scan potentially the entire database file. And if the database has not been compacted recently, there may be a large number of entries.

dbstats - the database statistics gathered by compact for the most used table columns

The contents of these tables may only be altered by the system.  However, they can be read from just like any other table.