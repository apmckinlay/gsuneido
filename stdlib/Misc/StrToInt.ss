// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (s)
	{
	return (s[0].Asc() << 24) + (s[1].Asc() << 16) + (s[2].Asc() << 8) + s[3].Asc()
	}