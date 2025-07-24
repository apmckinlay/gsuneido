// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// Contributed by Ajith.R
function (query, record)
	{
	return QueryDo('delete ' $ KeyQuery(query, record))
	}
