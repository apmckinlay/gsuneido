// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (book)
	{
	QueryApply(book $ " where not path.Prefix?('<deleted>')")
		{|x|
		if x.path is ""
			continue
		path = x.path.BeforeLast('/')
		if path is ""
			continue
		name = x.path.AfterLast('/')
		if QueryEmpty?(book, :path, :name)
			Print(x.path $ " (" $ x.name $ ")")
		}
	}