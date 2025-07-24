// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// WARNING: this only works if the query/view includes all records
// i.e. doesn't have a where or join
// SEE: IsDuplicateViaOutput for an alternative
function (query, field, value)
	{
	args = Object(query)
	args[field] = value
	return false is QueryEmpty?(@args)
	}
