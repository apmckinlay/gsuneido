// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (query, where)
	{
	// use \n in case query/where ends with //comment
	if BuiltDate() > #20250422
		return query $ '\n' $ where

	i = ScannerFind(query, "sort")
	return query[.. i] $ '\n' $ where $ '\n' $ query[i ..]
	}
