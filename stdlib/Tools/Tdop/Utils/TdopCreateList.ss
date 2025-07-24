// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	children = Object()
	block(children)
	return TdopCreateNode(TDOPTOKEN.LIST, :children)
	}
