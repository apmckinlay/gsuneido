// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	q = function (lib)
		{ return lib $ ' where group = -1' }
	n = 0
	libs = Object()
	for lib in Libraries()
		n += (libs[lib] = QueryCount(q(lib)))
	i = Random(n)
	for lib in libs.Members()
		if i < libs[lib]
			return QueryNth(i, q(lib)).Add(lib, at: #lib)
		else
			i -= libs[lib]
	}
