// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (query, field, value)
	{
	fn = Global(OptContribution('IsDuplicate', 'IsDuplicateViaQuery'))
	return fn(query, field, value)
	}
