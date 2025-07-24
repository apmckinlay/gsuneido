// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (query)
	{
	x = Query1(QueryStripSort(query) $ '\nsummarize count')
	return x is false ? 0 : x.count
	}