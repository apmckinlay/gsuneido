// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// similar to QueryAddWhere but handles more than where
// NOTE: use QueryAddWhere if just adding a where
function (query, more)
	{
	// use \n in case query/more ends with //comment
	return Query.StripSort(query) $
		Opt("\n", more) $
		Opt("\nsort ", Query.GetSort(query))
	}
