// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	RetryTransaction()
		{|t|
		return t.QueryDo(@args)
		}
	}