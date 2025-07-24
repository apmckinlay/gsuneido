// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// DEPRECATED
function (query, accum, block)
	{
	Transaction(read:)
		{ |t|
		t.QueryAccum(query, accum, block)
		}
	}

