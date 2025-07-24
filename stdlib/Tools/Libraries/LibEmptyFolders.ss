// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(lib)
		{
		.ForEach(lib)
			{
			Print(LibParentToPath(lib, it.parent) $ "/" $ it.name)
			}
		}
	ForEach(lib, block)
		{
		QueryApply(lib $ " where group >= 0") // folders
			{|x|
			if x.group is 'libcommitdate' or
				(not x.name.Suffix?("trial changes") and QueryEmpty?(lib, parent: x.num))
				block(x)
			}
		}
	RemoveAll()
		{
		n = 0
		for lib in LibraryTables()
			n += .Remove(lib)
		return n
		}
	Remove(lib)
		{
		n = 0
		.ForEach(lib)
			{
			Print("-", lib, LibParentToPath(lib, it.parent) $ "/" $ it.name)
			QueryDo("delete " $ lib $ " where num = " $ it.num)
			n++
			}
		return n
		}
	}