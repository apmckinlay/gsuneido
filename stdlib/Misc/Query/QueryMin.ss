// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (query, field, default = false)
	{
	Transaction(read:)
		{|t|
		t.QueryMin(query, field, default)
		}
	}