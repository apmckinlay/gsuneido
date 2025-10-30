<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Copy

``` suneido
() => record
```

Returns a copy of the record, including rule dependencies.

**Note:**

-	You can't call 
	[record.Update](<record.Update.md>) or 
	[record.Delete](<record.Delete.md>) on a copy of a record.
-	This does not copy observers.
-	This is a *shallow* copy i.e. if any members are themselves objects, they are not copied - they will just be references to the same objects.
-	[record.New?](<record.New?.md>) status is maintained


See also: [object.Copy](<../../../Language/Reference/Object/object.Copy.md>)