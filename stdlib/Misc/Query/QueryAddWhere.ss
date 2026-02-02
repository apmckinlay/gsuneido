// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (query, where)
	{
	// use \n in case query/where ends with //comment
	return query $ '\n' $ where
	}
