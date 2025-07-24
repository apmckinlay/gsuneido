// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (pattern)
	{
	node = Suneido.Parse('function () { ' $ pattern $ ' }')
	return node[0].type is 'MultiAssign' ? node[0] : node[0].expr
	}