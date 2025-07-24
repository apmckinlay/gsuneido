// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (query, record)
	{
	where = ShortestKey(query).
		Split(',').
		Map!({|f| f $ ' = ' $ Display(record[f]) }).
		Join(' and ')
	return where is ''
		? query
		: query $ ' where ' $ where
	}