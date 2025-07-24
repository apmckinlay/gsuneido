// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (lib, parent)
	{
	path = ""
	while parent isnt 0 and parent isnt 'libcommitdate'
		{
		x = Query1(lib, num: parent)
		path = "/" $ x.name $ path
		parent = x.parent
		}
	return path
	}