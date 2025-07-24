// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (x)
	{
	x &= 0xffff /*= 16 bits */
	// sign extend
	return (x ^ 0x8000) - 0x8000 /*= high bit */
	}
