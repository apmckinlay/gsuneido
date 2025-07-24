// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function (libRec, library)
	{
	path = libRec.name
	while false isnt libRec = Query1(library, num: libRec.parent)
		path = libRec.name $ '/' $ path
	return library $ '/' $ path
	}