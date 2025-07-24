// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (value)
	{
	return '0x' $ (value & 0xffffffff).Hex().LeftFill(8, '0') /*=
		need mask to be compatible with BuiltDate > 2025-05-09 */
	}
