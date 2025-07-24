// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
LibTreeModel
	{
	Children(parent)
		{ return super.Children(parent).Filter({ it.group }) }

	Children?(num)
		{
		lib = .TableName(num)
		if lib is false
			return false
		parent = .UnMangleNum(num)
		return not QueryEmpty?(lib $ ' where group > -1', :parent)
		}
	}