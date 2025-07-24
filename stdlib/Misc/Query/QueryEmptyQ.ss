// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20250331
function(@args)
	{
	if BuiltDate() < #20250424
		Transaction(read:)
			{ |t|
			return t.QueryEmpty?(@args)
			}
	return not QueryExists?(@args)
	}