// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	Transaction(read:)
		{ |t|
		return t.QueryAny1(@args)
		}
	}