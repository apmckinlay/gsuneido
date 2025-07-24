// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (x)
	{
	return Number?(x) or (String?(x) and x.Number?())
	}