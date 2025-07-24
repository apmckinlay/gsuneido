// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(ignoreDefault = false)
		{
		return .WithoutSpecial(Sys.Connections(), ignoreDefault)
		}
	WithoutSpecial(list, ignoreDefault = false)
		{
		ignoreList = Object("^[(].*[)]$", "^127.0.0.1", ':Thread-[0-9]+$')
		if ignoreDefault
			ignoreList.MergeUnion(GetContributions('UserConnections_IgnoreList'))
		return list.Copy().RemoveIf({|c| c =~ ignoreList.Join('|') or not c.Has?('@') })
		}
	}