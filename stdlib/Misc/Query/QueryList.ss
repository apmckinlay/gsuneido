// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// TODO: put field first and accept field: value shortcuts after query
function (query, field)
	{
	orig = query
	query = QueryStripSort(query)
	if query isnt orig
		Print('WARNING: QueryList ignores sort: ' $ orig)
	x = Query1(query $ '\nsummarize list ' $ field)
	return x is false ? [] : x['list_' $ field]
	}