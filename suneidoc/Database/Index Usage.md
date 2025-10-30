## Index Usage

An index consists of one or more columns in one table. A key is a kind of index.

The query optimizer will "choose" the index to use for a query based on the query, the conditions on the query, and the sort. The actual data also affects the choice of index. i.e. the strategy on one database may not be the same as the strategy on a different database.

Indexes are used for several purposes:

-	sort
-	where
-	operations like join or union


Indexes can be used in several ways:

-	to read records in a certain order (e.g. for sort)
-	to read a range of records (without reading the whole table)
-	to look up a specific record by a key


### where

Overall, **where** is looking for "x and y and z" where the parts like x are like \<field> \<op> \<constant>. i.e. it can use an index for "x > 1" but not "x > y" or "x > f()"

where can use indexes with:

-	>, >=, <, <=
-	"is" and "in" e.g. k in (1,2,3) where k is a key will do three lookups
-	"isnt" is treated as a range e.g. "a isnt 1" is treated as "a &lt 1 and a > 1"
-	"x is 1 or x is 2 or x is 3" is treated as x in (1,2,3)
-	[Date?](<../Language/Reference/Date?.md>), 
	[Number?](<../Language/Reference/Number?.md>), and 
	[String?](<../Language/Reference/String?.md>) are converted to ranges


The examples in [where_test.go](<https://github.com/apmckinlay/gsuneido/blob/master/dbms/query/where_test.go>) may be helpful.

### Query1, QueryEmpty?

-	using Query1/QueryEmpty? with just a table name in the first argument, and named arguments for the rest of the query is more efficient because it will bypass query parsing and optimizing
-	Query1 will output a slow query warning to error log if it reads more than 100 records to find what the query is looking for (as of 2025-09-12)
-	QueryEmpty? will output a slow query warning to error log if it read more than 2000 records to find what the query is looking for (as of 2025-09-12)


### Notes:

-	=~ and !~ will not use an index
-	an exact match key index will take priority (since it is guaranteed to not match more than one record)
-	you can use 
	[Trace](<../Language/Reference/Trace.md>)(TRACE.QUERYOPT) to get information on the optimizer's index estimates
-	In general, query execution can only use one index. If it uses a particular index to sort, then it can't use a different index for the where.
-	If you have index a,b you can read by "a" without having a separate index "a".  Indexes are read from left to right so the prefix of an index can be used without the whole index
-	Only the last field used in an index can be a range e.g. index(a,b,c,d) where a=1 and b=2 and c > 3 and d=4 can only use a,b,c and where a > 1 and ... can only use a
-	If you have an index(a,b) then `where a in (1,2) and b in (3,4)` will be "exploded" to the equivalent of: `a,b is 1,3 or a,b is 1,4 or a,b is 2,3 or a,b is 2,4`
-	You can use 
	[Trace](<../Language/Reference/Trace.md>)(TRACE.SLOWQUERIES) to see all slow queries
-	Using 
	[QueryStrategy](<Reference/QueryStrategy.md>) can help you figure out what index is being used BUT this can change depending on the actual data.  If you are using Query1, QueryEmpty?, QueryFirst, or QueryLast you can use 
	[Query.Strategy1](<Reference/Query/Query.Strategy1.md>) to see what index is being used.