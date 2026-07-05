// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	q = ' where group = -1'
	n = 0
	libs = Object()
	for lib in Libraries()
		n += (libs[lib] = QueryCount(lib $ q))
	i = Random(n)
	for lib in libs.Members()
		if i < libs[lib]
			return QueryNth(i, lib $ q).Add(lib, at: #lib)
		else
			i -= libs[lib]
	}
