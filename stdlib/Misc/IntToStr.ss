// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (n)
	{
	return (n >> 24).Chr() $ (n >> 16).Chr() $ (n >> 8).Chr() $ n.Chr()
	}
