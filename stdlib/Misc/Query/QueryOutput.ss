// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (query, record)
	{
	RetryTransaction()
		{ |t|
		t.QueryOutput(query, record)
		}
	}