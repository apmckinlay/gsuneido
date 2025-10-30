#### object.Extract

``` suneido
(name [, default]) => value
```

Combines [object.GetDefault](<object.GetDefault.md>) and [object.Delete](<object.Delete.md>)

-	if the member exists, its value is returned and the member is deleted
-	if the member doesn't exist
	
	-	if you supplied a default, it is returned
	-	otherwise, an exception is thrown
	



See also:
[NameArgs](<../NameArgs.md>),
[object.ListToMembers](<object.ListToMembers.md>),
[object.ListToNamed](<object.ListToNamed.md>),
[object.Project](<object.Project.md>),
[object.ProjectValues](<object.ProjectValues.md>),
object.Slice
