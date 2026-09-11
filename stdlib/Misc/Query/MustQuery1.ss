// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(@args)
	{
	x = Query1(@args)
	if x is false
		throw "MustQuery1: expected record, got false"
	return x
	}
