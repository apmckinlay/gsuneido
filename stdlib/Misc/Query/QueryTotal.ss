// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (query, field)
	{
	x = Query1(QueryStripSort(query) $ '\nsummarize total ' $ field)
	return x is false ? 0 : x['total_' $ field]
	}