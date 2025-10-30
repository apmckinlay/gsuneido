<div style="float:right"><span class="builtin">Builtin</span></div>

### Object

|     |     |     |
| --- | --- | --- |
| [Object](<Object/Object.md>) | [object.FindIf](<Object/object.FindIf.md>) | [object.NotEmpty?](<Object/object.NotEmpty?.md>) |
| [object.Add](<Object/object.Add.md>) | [object.FindLastIf](<Object/object.FindLastIf.md>) | [object.Nth](<Object/object.Nth.md>) |
| [object.AddMany!](<Object/object.AddMany!.md>) | [object.FindOne](<Object/object.FindOne.md>) | [object.PopFirst](<Object/object.PopFirst.md>) |
| [object.AddTo](<Object/object.AddTo.md>) | [object.First](<Object/object.First.md>) | [object.PopLast](<Object/object.PopLast.md>) |
| [object.AddUnique](<Object/object.AddUnique.md>) | [object.FlatMap](<Object/object.FlatMap.md>) | [object.Project](<Object/object.Project.md>) |
| [object.Any?](<Object/object.Any?.md>) | [object.Flatten](<Object/object.Flatten.md>) | [object.ProjectValues](<Object/object.ProjectValues.md>) |
| [object.Append](<Object/object.Append.md>) | [object.Fold](<Object/object.Fold.md>) | [object.RandVal](<Object/object.RandVal.md>) |
| [object.Assocs](<Object/object.Assocs.md>) | [object.GetDefault](<Object/object.GetDefault.md>) | [object.Readonly?](<Object/object.Readonly?.md>) |
| [object.BinarySearch](<Object/object.BinarySearch.md>) | [object.GetInit](<Object/object.GetInit.md>) | [object.Reduce](<Object/object.Reduce.md>) |
| [object.BinarySearch?](<Object/object.BinarySearch?.md>) | [object.Has?](<Object/object.Has?.md>) | [object.Remove](<Object/object.Remove.md>) |
| [object.CompareAndSet](<Object/object.CompareAndSet.md>) | [object.HasIf?](<Object/object.HasIf?.md>) | [object.RemoveIf](<Object/object.RemoveIf.md>) |
| [object.Concat](<Object/object.Concat.md>) | [object.HasNamed?](<Object/object.HasNamed?.md>) | [object.Replace](<Object/object.Replace.md>) |
| [object.Copy](<Object/object.Copy.md>) | [object.HasNonEmptyMember?](<Object/object.HasNonEmptyMember?.md>) | [object.Reverse!](<Object/object.Reverse!.md>) |
| [object.Count](<Object/object.Count.md>) | [object.Intersect](<Object/object.Intersect.md>) | [object.Set_default](<Object/object.Set_default.md>) |
| [object.CountIf](<Object/object.CountIf.md>) | [object.Intersects?](<Object/object.Intersects?.md>) | [object.Set_readonly](<Object/object.Set_readonly.md>) |
| [object.DeepCopy](<Object/object.DeepCopy.md>) | [object.Iter](<Object/object.Iter.md>) | [object.Shuffle!](<Object/object.Shuffle!.md>) |
| [object.Delete](<Object/object.Delete.md>) | [object.Join](<Object/object.Join.md>) | [object.Size](<Object/object.Size.md>) |
| [object.DeleteIf](<Object/object.DeleteIf.md>) | [object.JoinCSV](<Object/object.JoinCSV.md>) | [object.Sort!](<Object/object.Sort!.md>) |
| [object.Difference](<Object/object.Difference.md>) | [object.Last](<Object/object.Last.md>) | [object.SortWith!](<Object/object.SortWith!.md>) |
| [object.Disjoint?](<Object/object.Disjoint?.md>) | [object.ListToMembers](<Object/object.ListToMembers.md>) | [object.Sorted?](<Object/object.Sorted?.md>) |
| [object.Drop](<Object/object.Drop.md>) | [object.ListToNamed](<Object/object.ListToNamed.md>) | [object.Subset?](<Object/object.Subset?.md>) |
| [object.DuplicateValues](<Object/object.DuplicateValues.md>) | [object.Map](<Object/object.Map.md>) | [object.Sum](<Object/object.Sum.md>) |
| [object.Each](<Object/object.Each.md>) | [object.Map!](<Object/object.Map!.md>) | [object.SumWith](<Object/object.SumWith.md>) |
| [object.Each2](<Object/object.Each2.md>) | [object.Map2](<Object/object.Map2.md>) | [object.Swap](<Object/object.Swap.md>) |
| [object.Empty?](<Object/object.Empty?.md>) | [object.MapMembers](<Object/object.MapMembers.md>) | [object.Take](<Object/object.Take.md>) |
| [object.EqualSet?](<Object/object.EqualSet?.md>) | [object.Max](<Object/object.Max.md>) | [object.Trim!](<Object/object.Trim!.md>) |
| [object.Erase](<Object/object.Erase.md>) | [object.MaxWith](<Object/object.MaxWith.md>) | [object.Union](<Object/object.Union.md>) |
| [object.Eval](<Object/object.Eval.md>) | [object.Member?](<Object/object.Member?.md>) | [object.Unique!](<Object/object.Unique!.md>) |
| [object.Eval2](<Object/object.Eval2.md>) | [object.Members](<Object/object.Members.md>) | [object.UniqueValues](<Object/object.UniqueValues.md>) |
| [object.Every?](<Object/object.Every?.md>) | [object.MembersIf](<Object/object.MembersIf.md>) | [object.Val_or_func](<Object/object.Val_or_func.md>) |
| [object.Extract](<Object/object.Extract.md>) | [object.Merge](<Object/object.Merge.md>) | [object.Values](<Object/object.Values.md>) |
| [object.Filter](<Object/object.Filter.md>) | [object.MergeNew](<Object/object.MergeNew.md>) | [object.Without](<Object/object.Without.md>) |
| [object.Find](<Object/object.Find.md>) | [object.MergeUnion](<Object/object.MergeUnion.md>) | [object.WithoutFields](<Object/object.WithoutFields.md>) |
| [object.FindAll](<Object/object.FindAll.md>) | [object.Min](<Object/object.Min.md>) | [object.Zip](<Object/object.Zip.md>) |
| [object.FindAllIf](<Object/object.FindAllIf.md>) | [object.MinWith](<Object/object.MinWith.md>) |



Note: User defined methods for Object are defined in the stdlib **`Objects`** class.

Object constants are written with #(...)

Objects can be created with Object(...) or with [...]

[...] when it has **unnamed** members, is a shortcut for Object(...)

``` suneido
Type([1, 2, 3]) => "Object"
Type([1, 2, a: 3]) => "Object"
```

[...] with **only** named members (or no members) is a shortcut for Record(...)

``` suneido
Type([a: 1, b: 2]) => "Record"
```

**Note:** Named members are in a hash table and do not have any particular order.

See also: [Record](<../../Database/Reference/Record.md>)