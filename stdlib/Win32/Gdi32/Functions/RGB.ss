// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (r, g, b)
	{
	return (0xff & r) | ((0xff & g) << 8) | ((0xff & b) << 16)
	}
