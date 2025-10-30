### Syntax

*admin* = 
<pre>
<b>create</b> <i>table</i> <i>tablespec</i>
<b>ensure</b> <i>table</i> <i>tablespec</i>
<b>alter</b> <i>table</i> <b>create</b> <i>tablespec</i>
<b>alter</b> <i>table</i> <b>drop</b> <i>tablespec</i>
<b>alter</b> <i>table</i> <b>rename</b> <i>oldcolname</i> <b>to</b> <i>newcolname</i>
<b>view</b> <i>table</i> = <i>query</i>
<b>drop</b> <i>table</i>
<b>rename</b> <i>oldtablename</i> <b>to</b> <i>newtablename</i>
</pre>

*tablespec* =
<pre>
( <i>columns</i> )
<b>key</b> ( <i>columns</i> )
<b>index</b> [ <b>unique</b> ] ( <i>columns</i> ) [ <b>in</b> <i>table</i> [ ( <i>columns</i> ) ] ]
</pre>

-	Multiple keys and indexes may be specified.  Indexes are not a part of the "logical" design of the database.  Adding or removing indexes has no affect on the operation of the database other than on how fast certain queries can be executed.
-	Specifying unique on an index means that if the value is not empty, it must be unique, i.e. the only duplicates allowed are empty values.
-	Suneido's database does not require you to choose a "primary key"; but it requires at least one "candidate key" on each table, and it is a good idea to specify each candidate key.
-	Capitalized column names specify derived columns. These are calculated using a function called "Rule_" $ colname - they are not physically stored in the table. They are accessed as normal fields, without the capitalization. Rules are called as if they were methods of the record, i.e. *this *will be the record, so other members can be accessed as ".name".
-	If a column name ends with "_lower!" it is an automatically derived column. It is not physically stored in the database. Instead its value is the lower case version of the column with the same name without the "_lower!" (which must exist). Note: These fields may be indexed (unlike capitalized derived rule columns). For example:
	
	``` suneido
	create (name, name_lower!, age) key(name_lower!)
	```