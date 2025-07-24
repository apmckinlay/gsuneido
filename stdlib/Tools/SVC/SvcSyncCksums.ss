// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// Checksum goes: name, path, text (timmed with Trim()), lib committed date (string)
Memoize
	{
	Func(master)
		{
		names = Object()
		checksums = Object()
		book? = QueryColumns(master).Has?('svc_book')
		QueryApply(SvcCore.CurrentLibraryRecordsQuery(master))
			{|x|
			names.Add(x.name)
			if book?
				SvcBook.SplitText(x)
			checksums.Add(SvcCksum(x))
			}
		return [:names, :checksums]
		}
	}
