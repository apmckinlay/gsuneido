// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (library, x, t = false)
	{
	DoWithTran(t, update:)
		{|t|
		x = x.Copy()
		x.GetInit('num', { t.QueryMax(library, 'num', 0) + 1 })
		if x.Member?(#group)
			x.parent = x.group
		else
			x.group = -1
		x.GetInit('parent', 0)
		x.GetInit('lib_modified', { Date() })
		LibUnload(x.name)
		t.QueryOutput(library, x)
		}
	}
