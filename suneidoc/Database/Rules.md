## Rules

-	You can define *rules* for fields by defining functions called Rule_*fieldname* or by explicitly associating a rule with a field using 
	[record.AttachRule](<Reference/Record/record.AttachRule.md>)
-	When you access a field that the record does <u>not</u> contain, 
	if there is a rule it will be called.
-	If the rule returns a value, it will be stored in that field of the record.
-	When rules are executed, Suneido automatically tracks their *dependencies*on other fields they access. This includes any fields accessed by functions called by the rule.
-	Dependencies can also be set manually using 
	[record.SetDeps](<Reference/Record/record.SetDeps.md>) (and retrieved with 
	[record.GetDeps](<Reference/Record/record.GetDeps.md>))
-	If a dependency is changed, then the rule field is *invalidated*.
	This means that the next time the field is accessed, the rule will be executed again.
-	Dependencies can be stored in the database (by creating a field called *fieldname*_deps)
	so that when old records are manipulated,
	rules will be triggered just as on new records. 
	See: 
	[Stored Dependencies](<Rules/Stored Dependencies.md>)
-	Invalidations also trigger
	[record.Observer](<Reference/Record/record.Observer.md>)- this is used to update the user interface when records change.
-	Rules should not have *side effects* i.e. they should not modify any other fields.
-	Rules can be evaluated on record constants #{...} but will be slower because they cannot cache results.


Rules can be used without actually storing the values:
[Derived Columns](<Rules/Derived Columns.md>),
or calculated columns can be stored in the database:
[Stored Rules](<Rules/Stored Rules.md>)

Rules can also be used to adjust controls:
[Control Rules](<Rules/Control Rules.md>)