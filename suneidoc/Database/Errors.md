## Errors

The database functions throw exceptions if they encounter errors.
This means you don't have to always check return values.
The occasional time you do want to allow for errors, you can use try-catch or
[Catch](<../Language/Reference/Catch.md>).

**Some general errors:**

`nonexistent table:` *table*

>	Can't access a table that does not exist. 

`table already exists:` *table*

>	Can't rename to or create a table that already exists.

`nonexistent column:` *column*

>	Can't access a column that does not exist. 

`column already exists:` *column*

>	Can't rename to or create a column that already exists.

`duplicate key:` *keyfields* = *keyvalue*

>	Could not add a key or add or update a record due to a duplicate key value.

`blocked by foreign key:` *keyvalue*

>	The action could not be performed because it would not have satisfied the foreign key constraints.

`transaction conflict`

>	The action could not be performed because it would conflict with another transaction.
>	For example, deleting or updating a record that another transaction has deleted or updated.

`transaction aborted: transaction exceeded max age
update transaction longer than 20 seconds`

>	*As Of July 2022*   
>	The transaction was open too long

`too many overlapping update transactions`

>	Too many update transactions overlapping each other.  The maximum number of overlapping transactions is 500. This includes completed transactions if they overlap with an uncompleted transaction.

`too many reads/writes`

>	In an update transaction, the write limit is 10,000 records and the read limit is 20,000 records.  There are no limits on read-only transactions.

**Errors from requests:**

`query: create: key required for:` *table*

>	Tables must be created with at least one key.

`query: alter: rename column: can't rename system column:` *column*

`query: alter: add index: nonexistent column(s):` *columns*

>	The columns in an index or key must exist in the table.

`query: alter: add index: index already exists:` *index*

`query: alter: delete column: can't delete column used in index:` *column*

`query: alter: delete column: can't delete system column:` *column*

`query: alter: delete index: nonexistent index:` *index*

`query: alter: delete index: can't delete system index:` *index*

`query: alter: delete index: can't delete last index:` *index*

`query: rename table: can't rename system table:` *table*

`query: drop: can't destroy system table:` *table*

**Errors from queries:**

`query: rename: column(s) already exist:` *columns*

`query: extend: column(s) already exist:` *columns*

`query: join: common columns required`

`query: product: common columns not allowed:` *columns*

`
query: project: output: key required
query: project: update: key required
query: project: delete: key required`

>	You can only modify through a project if the columns include a key.

`difference: inputs must have the same columns`

`union: inputs must have the same columns`

`intersect: inputs must have the same columns`

`can't add records to system table:` *table*   
`can't update records in system table:` *table*   
`can't delete records from system table:` *table*

>	System tables are maintained automatically.

`
output is not allowed on this query
update is not allowed on this query
delete is not allowed on this query
`

>	Not all queries support modifications.

**temp index: derived too large:**

"derived" is data that is not from the database, usually from extend although potentially also from summarize.

Temp indexes normally only store and sort pointers to the data in the database. But derived data isn't in the database so it has to store the data itself in memory.

Temp indexes often come from a sort on the query, but they can also be necessary for intermediate steps in complex queries.

Because there can be many users (potentially hundreds) and because the server has a finite amount of memory, Suneido has to limit how much memory is uses. Otherwise, the worst case is that the server runs out of memory and crashes. Even if it doesn't crash, excessive memory use will slow down the server for everyone.

One option when you encounter this error is to use a select or summarize to limit the number of records included in the query.

Another option is to eliminate or reduce the extends. In some cases you may be able to replace them with rules. Rules are calculated on demand rather than being stored, so they are not included in derived. However, rules will likely be slower, which is not good since sorting huge amounts of data is already going to be slow.

In some cases you can add an index to a table so the temp index is not required. (Although index usage also depends on other factors like selects.)

*As Of September 2023*

temp index warning at 200,000 and temp index derived warning at 8,000,000