// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (query)
	{
	result = QueryAddKeySort(query)
	if result is false
		throw 'QueryEnsureKeySort: query has non-key sort'
	return result
	}