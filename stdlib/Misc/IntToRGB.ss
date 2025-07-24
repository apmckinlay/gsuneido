// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (int)
	{
	base = 256
	r = int % base
	g = (int / base).Int() % base
	b = (int / base / base).Int() % base
	return Object(:r, :g, :b)
	}