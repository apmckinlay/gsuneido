// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (query, field)
	{
	Transaction(read:)
		{|t|
		t.QueryAverage(query, field)
		}
	}