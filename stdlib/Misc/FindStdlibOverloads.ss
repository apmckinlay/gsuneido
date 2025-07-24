// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (refPrev? = false)
	{
	stdlib = Object()
	QueryApply('stdlib', group: -1)
		{|x|
		stdlib[x.name] = Object()
		}
	for lib in LibraryTables().Remove('stdlib')
		QueryApply(lib, group: -1)
			{|x|
			if stdlib.Member?(x.name) and
				(not refPrev? or x.lib_current_text.Has?('_' $ x.name))
				stdlib[x.name].Add(lib)
			}
	for name in stdlib.Members().Sort!()
		if stdlib[name].NotEmpty?()
			Print(name.RightFill(30/*= max name*/), stdlib[name].Join(', '))
	}
