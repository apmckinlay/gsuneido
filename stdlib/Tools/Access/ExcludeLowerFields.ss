// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
function (columns, excludeFields = #())
	{
	return columns.Filter({ it.Suffix?("_lower!") }).Add(@excludeFields)
	}
