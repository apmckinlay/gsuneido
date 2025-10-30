### Suneido

An object that can be used to store "global" information.

Since this object is shared, make sure member names are unique so you don't overwrite other information. e.g. prefix member names with the name of your library record.

If you need to store several items, consider grouping them into an object so you are only adding one member to Suneido.

**Note**: Global variables are often a bad idea, consider alternatives.

#### Methods

|     |     |     |
| --- | --- | --- |
| [Suneido.AssertFail](<Suneido/Suneido.AssertFail.md>) | [Suneido.LibraryTags](<Suneido/Suneido.LibraryTags.md>) | [Suneido.ShouldNotReachHere](<Suneido/Suneido.ShouldNotReachHere.md>) |
| [Suneido.Crash!](<Suneido/Suneido.Crash!.md>) | [Suneido.Parse](<Suneido/Suneido.Parse.md>) | [Suneido.StrictCompare](<Suneido/Suneido.StrictCompare.md>) |
| [Suneido.GoMetric](<Suneido/Suneido.GoMetric.md>) | [Suneido.Regex](<Suneido/Suneido.Regex.md>) | [Suneido.WarningsThrow](<Suneido/Suneido.WarningsThrow.md>) |
| [Suneido.Info](<Suneido/Suneido.Info.md>) | [Suneido.RuntimeError](<Suneido/Suneido.RuntimeError.md>) |

