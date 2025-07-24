// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (lo, hi)
	{
	mask = 0xffff
	return ((hi & mask) << 16) | (lo & mask)
	}
