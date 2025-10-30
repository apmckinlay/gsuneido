<div style="float:right"><span class="builtin">Builtin</span></div>

### Record

|     |     |
| --- | --- |
| [Record](<Record/Record.md>) | [record.Observer](<Record/record.Observer.md>) |
| [record.AttachRule](<Record/record.AttachRule.md>) | [record.PreSet](<Record/record.PreSet.md>) |
| [record.Clear](<Record/record.Clear.md>) | [record.Project](<Record/record.Project.md>) |
| [record.Copy](<Record/record.Copy.md>) | [record.RemoveObserver](<Record/record.RemoveObserver.md>) |
| [record.Delete](<Record/record.Delete.md>) | [record.SafeMembers](<Record/record.SafeMembers.md>) |
| [record.DoWithTran](<Record/record.DoWithTran.md>) | [record.SetDeps](<Record/record.SetDeps.md>) |
| [record.GetDeps](<Record/record.GetDeps.md>) | [record.Transaction](<Record/record.Transaction.md>) |
| [record.Invalidate](<Record/record.Invalidate.md>) | [record.Update](<Record/record.Update.md>) |
| [record.New?](<Record/record.New?.md>) |



Record is derived (inherits) from [Object](<../../Language/Reference/Object.md>) and has all the Object methods plus some extra features:

-	default value of "" equivalent to 
	[object.Set_default("")](<../../Language/Reference/Object/object.Set_default.md>)
-	observers
-	rules and dependencies
-	database facilities


It is reasonable to use Record's for their observer and rule facilities, independent of the database. (e.g. RecordControl) However, there is some performance cost for these extra features so you would normally not use Records unless you need these features.

Record constants are written as #{...} 

Records can be created with Record(...) or with [...]

[...] when it has **unnamed** members, is a shortcut for Object(...)

``` suneido
Type([1, 2, 3]) => "Object"
Type([1, 2, a: 3]) => "Object"
```

[...] with **only** named members (or no members) is a shortcut for Record(...)

``` suneido
Type([a: 1, b: 2]) => "Record"
```

User defined methods for Record may be added in the **`Records`** class.

See also: [Object](<../../Language/Reference/Object.md>)