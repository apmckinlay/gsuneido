## Introduction

Suneido has an integrated client-server relational database. The database is accessed via a language that includes administration requests and update operations as well as queries. The query language is based on the relational algebra language suggested by C.J.Date.

As a relational database, a Suneido database holds all of its data in tables (also known as relations). The actual data is contained in the rows or records of the tables (also known as tuples). A table has a set of columns or fields (also known as attributes). Tables must have at least one *key* consisting of one or more columns that uniquely identify records.

Columns (fields) can store the following types of values:

-	boolean (true or false)
-	string (including binary, e.g. image data)
-	number
-	date/time
-	object (i.e. arrays or records)


Suneido's DBMS, like its language, is dynamically typed i.e. database columns (fields) do not have fixed types - they can hold any type of value. Validating data is left up to the application. All fields and records are variable length.

Suneido stores the entire database as a single operating system file. This includes the schema (table layout) information, data records, indexes, and concurrency & recovery information. The database is accessed as a memory mapped file.

The database can operate in single-user local mode, or in multi-user client-server mode. TCP/IP is used to communicate between the clients and the server.

NOTE: Although objects can be stored in the database, they will not sort properly (as they do in the language). Indexes and sorting order values based on their packed format (see [Pack](<../Language/Reference/Pack.md>)). This works for simple values but not objects.

### Query Optimization

Query optimization has two main phases. The first phase applies some standard transformations to the query that are almost always advantageous. For example, moving where's towards tables and combining adjacent operations. In the second phase, operations choose appropriate strategies, indexes, and temporary indexes based on estimated costs. Data sizes are estimated using the indexes. Some operations have multiple strategies (e.g. project) they can use, other operations (e.g. rename) have only a single strategy.

### Rules

Unlike many systems, which limit business rules to constraints, Suneido's business rules support a variety of uses including supplying default values to fields, performing calculations, and summarizing other data. Business rules have many advantages. They keep your business logic separate from your user interface and reports, enable code re-use, and allow your code to be written in smaller modules that are easier to test and maintain.

You can define rules for fields by defining functions called Rule_fieldname. When you access a field that the record does not contain, if there is a rule it will be called. If the rule returns a value, it will be stored in that field of the record. When rules are executed, Suneido automatically tracks their dependencies on other fields they access. If a dependency is changed, then the rule field is invalidated. This means that the next time the field is accessed, the rule will be executed again. Dependencies can be stored in the database (by creating a field called fieldname_deps) so that when old records are manipulated, rules will be triggered just as on new records. Invalidations also trigger record.Observer - this is used to update the user interface when records change. Invalidations do not affect non-rule values. i.e. if the user has overridden a derived value, then the rule on that field will no longer be triggered. Rules can be used without actually storing the values, or calculated columns can be stored in the database. Rules can also be used to adjust user interface controls.

### Triggers

Whenever records are output, updated, or deleted from a table, the system checks for a user defined trigger named "Trigger_" followed by the table name. Triggers are only called after the operation succeeds. i.e. If the output, update, or delete fails, the trigger will not be called. Triggers can be used to maintain secondary tables such as summaries.

### History

[transaction.Asof](<Reference/Transaction/transaction.Asof.md>) allows read-only transactions as of some previous point in time. The available points in time are given by the "history" virtual table.

**Note**: Compacting the database removes the history.

### Concurrency

Suneido's DBMS can operate in one of two modes - single-user standalone mode, or multi-user client-server mode. In either case, the database file itself is only ever accessed by a single program exclusively.

All access to the database must be done within transactions. Transactions can be either read-only or update. Transactions see the database as of their start time, as if they were viewing a "snapshot" of the database. Suneido uses multi-version optimistic concurrency, which provides full transaction isolation, i.e. is serializable. Because of this, read-only transactions (e.g. reports) always succeed - they will never conflict with other transactions. Update transactions check for conflicts with other transactions, and fail (rollback) if conflicts are found.

On-line backups are done using a single read-only transaction to get a "snap-shot" of the entire database without interfering with use of the database.

### References

An Introduction to Database Systems, C.J.Date